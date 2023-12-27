package structuredLogging

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/itdesign-at/golib/keyvalue"
)

type SlogLogger struct {
	// writers holds all configured writer functions
	writers map[string]func(string, []byte)
	params  keyvalue.Record
}

// New creates a writer slog
// Examples:
//
//	structuredLogging.New("STDERR")
//	structuredLogging.New("STDOUT")
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
		params:  make(keyvalue.Record),
	}

	for _, d := range dsn {
		// additional to slices, comma seperated strings are supported
		for _, str := range strings.Split(d, ",") {
			switch {
			case strings.HasPrefix(str, "/"):
				// it's a file
				sl.writers[str] = sl.writeFile
			case strings.HasPrefix(str, "nats://"):
				// it's nats
				sl.writers[str] = sl.writeNats
			case strings.ToUpper(str) == "STDOUT":
				// whatever, definitely stderr
				sl.writers["STDOUT"] = sl.writeStdOut
			default:
				// whatever, definitely stderr
				sl.writers["STDERR"] = sl.writeStdErr
			}
		}
	}

	return &sl
}

// Parameter sets structuredLogging parameters. Keys are converted to lowercase.
// Key/value pairs implemented:
//
//	"level": "debug" ... which level to log. "debug" is default
//	"level": "error"
//	"level": "warning"
//	"level": "info"
//	"natsSubject": "string" ... alternative NATS subject
//
// Examples:
//
//	  structuredLogging.New("STDERR").Parameter(
//			map[string]interface{}{"level": "error"}).Init()
//
//	  structuredLogging.New("nats://127.0.0.1").Parameter(
//			map[string]interface{}{"natsSubject": "messages.watchit"}).Init()
func (sl *SlogLogger) Parameter(params keyvalue.Record) *SlogLogger {
	for k, v := range params {
		sl.params[strings.ToLower(k)] = v
	}
	return sl
}

// Init initialize the logger
// In its simplest form structuredLogging.New().Init() logs to STDERR
func (sl *SlogLogger) Init() *SlogLogger {

	// if no writer is defined, it is ensured that logging occurs on stderr
	if len(sl.writers) == 0 {
		sl.writers["STDERR"] = sl.writeStdErr
	}

	var level slog.Level
	switch sl.params.String("level", true) {
	case "debug", slog.LevelDebug.String():
		level = slog.LevelDebug
	case "error", slog.LevelError.String():
		level = slog.LevelError
	case "info", slog.LevelInfo.String():
		level = slog.LevelInfo
	case "warning", slog.LevelWarn.String():
		level = slog.LevelWarn
	default:
		level = slog.LevelDebug
	}

	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
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
func (sl *SlogLogger) writeNats(dsn string, b []byte) {

	conn, err := nats.Connect(dsn, nats.RetryOnFailedConnect(false),
		nats.MaxReconnects(3))

	if err != nil {
		return
	}

	var subject string

	if subject = sl.params.String("natssubject", true); subject == "" {
		var x map[string]any
		_ = json.Unmarshal(b, &x)
		subject = "slog."
		if level, ok := x["level"].(string); ok {
			subject += level
		} else {
			subject += "UNKNOWN"
		}
	}

	_ = conn.Publish(subject, b)
	conn.Close()

}

// writeFile write buffer b to file f
func (sl *SlogLogger) writeFile(f string, b []byte) {
	if file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		_, _ = file.Write(b)
		_ = file.Close()
	}
}

// writeStdErr write buffer b to os.Stderr
func (sl *SlogLogger) writeStdErr(f string, b []byte) {
	_, _ = os.Stderr.Write(b)
	// NEVER close os.Stderr, see package os
}

// writeStdErr write buffer b to os.Stdout
func (sl *SlogLogger) writeStdOut(f string, b []byte) {
	_, _ = os.Stdout.Write(b)
	// NEVER close os.Stdout, see package os
}
