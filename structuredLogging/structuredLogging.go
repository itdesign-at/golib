package structuredLogging

import (
	"encoding/json"
	"log/slog"
	"math"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type SlogLogger struct {
	// writers holds all configured writer functions
	writers map[string]func(string, []byte)
}

// New creates a writer slog
// Examples:
//
//	structuredLogging.New("STDERR")
//	structuredLogging.New("/var/log/myLogfile.log")
//	structuredLogging.New("nats://server1.demo.at:4222/subject.prefix")
//	structuredLogging.New([]string{"STDERR","/var/log/myLogFile.log"}...)
//	structuredLogging.New("STDERR","/var/log/myLogFile.log","/var/log/anotherLogFile.log")
//	structuredLogging.New("STDERR","/var/log/myLogFile.log")
//
// In its simplest form structuredLogging.New().Init() logs to STDERR
func New(dsn ...string) *SlogLogger {
	var sl = SlogLogger{
		writers: make(map[string]func(string, []byte)),
	}

	for _, d := range dsn {
		// additional to slices, comma seperated strings are supported
		for _, str := range strings.Split(d, ",") {
			switch {
			case strings.HasPrefix(str, "/"):
				// it's a file
				sl.writers[str] = writeFile
			case strings.HasPrefix(str, "nats://"):
				// it's nats
				sl.writers[str] = writeNats
			default:
				// whatever, definitely stderr
				sl.writers["STDERR"] = writeStdErr
			}
		}
	}

	return &sl
}

// Parameter sets structuredLogging parameters
// under construction: currently no function implemented, reserved fpr future use
func (sl *SlogLogger) Parameter(params ...string) *SlogLogger {

	return sl
}

// Init initialize the logger
// In its simplest form structuredLogging.New().Init() logs to STDERR
func (sl *SlogLogger) Init() *SlogLogger {

	// if no writer is defined, it is ensured that logging occurs on stderr
	if len(sl.writers) == 0 {
		sl.writers["STDERR"] = writeStdErr
	}

	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}

	jh := slog.NewJSONHandler(sl, opt)
	logger := slog.New(jh)
	slog.SetDefault(logger)

	return sl
}

// GenerateLogfileName generates a filename with a golang
// time layout as input. "user.Current" is replaced with the username
// of the user.Current() method.
//
// Example:
//
//	/var/log/messenger-user.Current-2006-01.log
//
// returns
//
//	/var/log/messenger-root-2023-12.log
func GenerateLogfileName(layout string) string {
	var str string
	var username string

	if u, err := user.Current(); err == nil {
		username = u.Username
		if username == "" {
			username = u.Name
		}
		if username == "" {
			username = u.Uid
		}
	} else {
		uid := os.Getuid()
		username = strconv.Itoa(uid)
	}

	layout = strings.Replace(layout, "user.Current", username, -1)
	str = time.Now().Format(layout)

	return str
}

// Write to all configured log functions (sl.writers)
func (sl *SlogLogger) Write(p []byte) (int, error) {
	var wg sync.WaitGroup

	for dsn, f := range sl.writers {
		wg.Add(1)

		go func(f func(string, []byte), dsn string, p []byte) {
			f(dsn, p)
			wg.Done()
		}(f, dsn, p)
	}

	wg.Wait()
	return len(p), nil
}

// writeNats write buffer b to nats server (parameter dsn)
// the subject is extract from key "level" of the json message (buffer b)
// nats reconnect documentation see
// https://docs.nats.io/using-nats/developer/connecting/reconnect
func writeNats(dsn string, b []byte) {
	if conn, err := nats.Connect(dsn, nats.RetryOnFailedConnect(true), nats.MaxReconnects(math.MaxInt)); err == nil {
		var subject string
		var x map[string]any

		_ = json.Unmarshal(b, &x)
		if level, ok := x["level"].(string); ok {
			subject = "slog." + level
		} else {
			subject = "slog.UNKNOWN"
		}
		_ = conn.Publish(subject, b)
		conn.Close()
	}
}

// writeFile write buffer b to file f
func writeFile(f string, b []byte) {
	if file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		_, _ = file.Write(b)
		_ = file.Close()
	}
}

// writeStdErr write buffer b to os.Stderr
func writeStdErr(f string, b []byte) {
	_, _ = os.Stderr.Write(b)
	// NEVER close os.Stderr, see package os
}
