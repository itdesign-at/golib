package structuredLogging

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

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
func New(config map[string]any) *SlogLogger {
	var sl SlogLogger
	opt := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}
	if config == nil {
		jh := slog.NewJSONHandler(os.Stderr, opt)
		logger := slog.New(jh)
		slog.SetDefault(logger)
		return &sl
	}
	return &sl
}

func (sl SlogLogger) Write(p []byte) (int, error) {
	var data = make([]byte, len(p))
	copy(data, p)
	if sl.natsHandler != nil {
		go func(b []byte) {
			var subject string
			var x map[string]any
			_ = json.Unmarshal(b, &x)
			if level, ok := x["level"].(string); ok {
				subject = "slog." + level
			} else {
				subject = "slog.UNKNOWN"
			}
			_ = sl.natsHandler.Publish(subject, b)
		}(data)
	}
	if sl.fileHandler != nil {
		go func(b []byte) {
			_, _ = sl.fileHandler.Write(b)
		}(data)
	}
	if sl.stderrHandler != nil {
		_, _ = sl.stderrHandler.Write(data)
	}
	return len(p), nil
}
