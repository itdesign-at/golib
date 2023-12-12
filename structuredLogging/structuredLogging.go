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
	myConfig map[string]string

	w             io.Writer
	natsHandler   *nats.Conn
	fileHandler   *os.File
	stderrHandler *os.File
	message       chan []byte
}

// New takes a config string or slice with its parameters
// to init slog logging in json format.
//
// Examples:
//
//	structuredLogging.New("STDERR")
//	structuredLogging.New("/var/log/myLogfile.log")
//	structuredLogging.New("nats://server1.demo.at:4222/subject.prefix")
//
//	structuredLogging.New("STDERR,file:/var/log/myLogFile.log")
//
//	structuredLogging.New([]string{"STDERR","/var/log/myLogFile.log"})
//
// In its simplest form structuredLogging.New(nil) logs to STDERR
func New(config any) *SlogLogger {
	var sl SlogLogger

	sl.myConfig = make(map[string]string)

	var fillMyConfig = func(x string) {
		if strings.HasPrefix(x, "/") {
			sl.myConfig["file"] = x
		} else if strings.HasPrefix(x, "nats://") {
			sl.myConfig["nats"] = x
		} else {
			sl.myConfig["STDERR"] = "yes"
		}
	}

	switch x := config.(type) {
	case string:
		if strings.Contains(x, ",") {
			for _, str := range strings.Split(x, ",") {
				fillMyConfig(str)
			}
		} else {
			fillMyConfig(x)
		}
	case []string:
		for _, str := range x {
			fillMyConfig(str)
		}
	case nil:
		fillMyConfig("STDERR")
	}

	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}

	if n, ok := sl.myConfig["nats"]; ok && n != "" {
		// nats reconnect documentation see
		// https://docs.nats.io/using-nats/developer/connecting/reconnect
		sl.natsHandler, _ = nats.Connect(
			n,
			nats.RetryOnFailedConnect(true),
			nats.MaxReconnects(math.MaxInt),
		)
	}
	if f, ok := sl.myConfig["file"]; ok && f != "" {
		sl.fileHandler, _ = os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if _, ok := sl.myConfig["STDERR"]; ok {
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
