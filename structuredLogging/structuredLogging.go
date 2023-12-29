package structuredLogging

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/url"
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
	writers        map[string]func(string, []byte)
	params         keyvalue.Record
	natsConnection *nats.Conn
}

// New creates an slog writer which allows to log to multiple
// destinations.
//
// Examples:
//
//	structuredLogging.New("STDERR")
//	structuredLogging.New("STDOUT")
//	structuredLogging.New("/var/log/myLogfile.log")
//
//	structuredLogging.New("nats://server1.demo.at")
//	structuredLogging.New("nats://server1.demo.at/messenger.LOGLEVEL")
//	structuredLogging.New("nats://server1.demo.at/scheduler.demo")
//
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
		// comma seperated strings are supported
		for _, str := range strings.Split(d, ",") {
			switch {
			case strings.HasPrefix(str, "/"):
				sl.writers[str] = sl.writeFile
			case strings.HasPrefix(str, "nats://"):
				parse, err := url.Parse(str)
				if err != nil {
					sl.writers["STDOUT"] = sl.writeStdOut
					break
				}
				natsDSN := fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)
				if parse.Path != "" {
					sl.Parameter(map[string]interface{}{
						// e.g. parse.Path = "/mySubject.LOGLEVEL"
						"natssubject": strings.TrimPrefix(parse.Path, "/"),
					})
				}
				sl.writers[str] = sl.writeNats
				// nats reconnect documentation see
				// https://docs.nats.io/using-nats/developer/connecting/reconnect
				sl.natsConnection, _ = nats.Connect(natsDSN, nats.RetryOnFailedConnect(true),
					nats.MaxReconnects(math.MaxInt))
			case strings.ToUpper(str) == "STDOUT":
				sl.writers["STDOUT"] = sl.writeStdOut
			default:
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
//
//	  structuredLogging.New("nats://127.0.0.1").Parameter(
//			map[string]interface{}{"NatsSubject": "scheduler.LOGLEVEL"}).Init()
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
//	layout = /var/log/messenger-user.Current-2006-01.log
//
// returns
//
//	/var/log/messenger-root-2023-12.log
//
// With New() and Init():
//
//	logFilename := filepath.Join("/var/log",
//	  structuredLogging.GenerateLogfileName("messenger-user.Current-2006-01.log")
//	)
//	structuredLogging.New(logFilename).Init()
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

// writeNats write buffer jsonLogMessage to nats server which
// is connected above in the New() constructor.
// The subject is extracted from the "level" of the json message
// or derived from the "natssubject" parameter
func (sl *SlogLogger) writeNats(dsn string, jsonLogMessage []byte) {
	// dsn is the NATS destination which is currently unused
	subject := sl.prepareNatsSubject(jsonLogMessage)
	_ = sl.natsConnection.Publish(subject, jsonLogMessage)
}

// prepareNatsSubject processes the parameter "natssubject" and generates
// a subject. Default = "slog.<LOGLEVEL>" like e.g. slog.DEBUG
func (sl *SlogLogger) prepareNatsSubject(jsonLogMessage []byte) string {
	var subject = sl.params.String("natssubject", true)
	var extractLogLevel = func() string {
		var x map[string]any
		_ = json.Unmarshal(jsonLogMessage, &x)
		logLevel, _ := x["level"].(string)
		if logLevel == "" {
			return "UNKNOWN"
		}
		return logLevel
	}
	if subject == "" {
		subject = "slog." + extractLogLevel() // e.g. slog.INFO
	} else if strings.Contains(subject, "LOGLEVEL") {
		subject = strings.Replace(subject, "LOGLEVEL", extractLogLevel(), -1)
	}
	return subject
}

// writeFile write buffer "jsonLogMessage" to file "fileName"
func (sl *SlogLogger) writeFile(fileName string, jsonLogMessage []byte) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		_, _ = file.Write(jsonLogMessage)
		_ = file.Close()
	}
}

// writeStdErr write buffer "jsonLogMessage" to os.Stderr
func (sl *SlogLogger) writeStdErr(_ string, jsonLogMessage []byte) {
	_, _ = os.Stderr.Write(jsonLogMessage)
	// NEVER close os.Stderr, see package os
}

// writeStdErr write buffer "jsonLogMessage" to os.Stdout
func (sl *SlogLogger) writeStdOut(_ string, jsonLogMessage []byte) {
	_, _ = os.Stdout.Write(jsonLogMessage)
	// NEVER close os.Stdout, see package os
}
