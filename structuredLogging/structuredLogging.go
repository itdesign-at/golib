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

// TODO: Initialize() Methode
// TODO: Kommentare fehlen noch!

type SlogLogger struct {
	writers map[string]func(string, []byte)
}

// New takes a config string or slice with its parameters
// to init slog logging in json format.
//
// Examples:
//
//	structuredLogging.New("STDERR")
//	structuredLogging.New("/var/log/myLogfile.log")
//	structuredLogging.New("nats://server1.demo.at:4222/subject.prefix")
//	structuredLogging.New([]string{"STDERR","/var/log/myLogFile.log"}...)
//	structuredLogging.New("STDERR","/var/log/myLogFile.log","/var/log/anotherLogFile.log")
//	structuredLogging.New("STDERR","/var/log/myLogFile.log")
//	structuredLogging.New()
//
// In its simplest form structuredLogging.New() logs to STDERR
func New(dsn ...string) *SlogLogger {
	sl := SlogLogger{
		writers: make(map[string]func(string, []byte)),
	}

	var fillMyConfig = func(x string) {
		if strings.HasPrefix(x, "/") {
			// it's a file
			sl.writers[x] = writeFile
		} else if strings.HasPrefix(x, "nats://") {
			// it's nats
			sl.writers[x] = writeNats
		} else {
			// whatever, definitely stderr
			sl.writers["STDERR"] = writeStdErr
		}
	}

	for _, d := range dsn {
		if strings.Contains(d, ",") {
			// additional to slices, comma seperated strings are supported
			for _, str := range strings.Split(d, ",") {
				fillMyConfig(str)
			}
		} else {
			fillMyConfig(d)
		}
	}

	// if no writer is defined, it is ensured that logging occurs on stderr
	if len(sl.writers) == 0 {
		fillMyConfig("STDERR")
	}

	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}

	jh := slog.NewJSONHandler(sl, opt)
	logger := slog.New(jh)
	slog.SetDefault(logger)

	return &sl
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
	u, err := user.Current()
	if err == nil {
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

// Write to all configured log handlers
func (sl SlogLogger) Write(p []byte) (int, error) {
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

func writeFile(f string, b []byte) {
	if file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		_, _ = file.Write(b)
		_ = file.Close()
	}
}

func writeStdErr(f string, b []byte) {
	_, _ = os.Stderr.Write(b)
}
