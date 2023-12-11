package structuredLogging

import (
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type SlogLogger struct {
	w             io.Writer
	natsHandler   *nats.Conn
	fileHandler   *os.File
	stderrHandler *os.File
	message       chan []byte
}

// New takes a config map with its parameters to start
// Example:
//
//	structuredLogging.New(map[string]any{
//	    "nats":   "nats://127.0.0.1:4222",
//	    "file":   "/tmp/do-log.log",
//	    "STDERR": true,
//	})
//
// In its simplest form structuredLogging.New(nil) logs to STDERR
func New(config map[string]any) *SlogLogger {
	var sl SlogLogger

	if config == nil {
		config = map[string]any{"STDERR": true}
	}
	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}

	if n, ok := config["nats"].(string); ok && n != "" {
		// nats reconnect documentation see
		// https://docs.nats.io/using-nats/developer/connecting/reconnect
		sl.natsHandler, _ = nats.Connect(
			n,
			nats.RetryOnFailedConnect(true),
			nats.MaxReconnects(math.MaxInt),
		)
	}
	if f, ok := config["file"].(string); ok && f != "" {
		sl.fileHandler, _ = os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if b, ok := config["STDERR"].(bool); ok && b {
		sl.stderrHandler = os.Stderr
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
	var data = make([]byte, len(p))
	copy(data, p)
	if sl.natsHandler != nil {
		var subject string
		var x map[string]any
		_ = json.Unmarshal(data, &x)
		if level, ok := x["level"].(string); ok {
			subject = "slog." + level
		} else {
			subject = "slog.UNKNOWN"
		}
		_ = sl.natsHandler.Publish(subject, data)
	}
	if sl.fileHandler != nil {
		_, _ = sl.fileHandler.Write(data)
	}
	if sl.stderrHandler != nil {
		_, _ = sl.stderrHandler.Write(data)
	}
	return len(p), nil
}
