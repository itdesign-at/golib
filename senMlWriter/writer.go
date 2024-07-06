package senMlWriter

import (
	"errors"
	"fmt"
	"log/slog"
	"log/syslog"
	"strings"

	"github.com/mainflux/senml"
)

var ErrUnknownOut = errors.New("unknown Out format")

// WriterConfig holds the configuration for the Sensorml handler.
// The configuration is used to create a new Sensorml handler.
// The handler is used to add data to the Sensorml handler.
// The handler writes the data to the configured Out every minute.
// The Out should be in the form of "file:/tmp/astrolab%02d.json" or "syslog://localhost:5514/tag".
// The base name is used as the base name for the senml records.
// The debug flag is used to enable debug logging.
// The args are used to pass additional key value pairs to the handler.
// The args are not used by the handler.
type WriterConfig struct {
	// field(s) for senML
	BaseName string `json:"baseName" yaml:"baseName"`

	// field(s) for sending data
	Out string `json:"out" yaml:"out"`

	SyslogPriority syslog.Priority `json:"syslogPriority" yaml:"syslogPriority"`

	debug bool
}

type Writer struct {
	cfg WriterConfig
	p   senml.Pack
}

// NewWriter creates a new Writer with the given configuration and current time.
// The writer is used to write senml.Pack to the configured Out.
// The Out should be in the form of "file:/tmp/astrolab%02d.json" or "syslog://localhost:5514/tag".
// The current time is used to replace any time fields in the Out.
// The writer is not thread safe.
// Use a new writer for each thread.
func NewWriter(cfg WriterConfig) *Writer {
	if cfg.SyslogPriority == 0 {
		cfg.SyslogPriority = syslog.LOG_INFO | syslog.LOG_LOCAL7
	}
	return &Writer{
		cfg: cfg,
		p:   senml.Pack{},
	}
}

// AddPack adds a senml.Pack to the writer.
// The pack is written to the configured Out when Write is called.
// The pack is written as is, no sorting is done.
func (w *Writer) AddPack(p senml.Pack) *Writer {
	w.p = p
	if w.cfg.debug {
		slog.Debug("senMlWriter AddPack", "records", p.Len())
	}
	return w
}

// Write writes the senml.Pack to the configured Out.
// It returns an error if the writing fails.
// The Out should be in the form of "file:/tmp/astrolab%02d.json" or "syslog://localhost:5514/tag".
func (w *Writer) Write() error {
	var err error
	// w.cfg.Out examples:
	// file:/tmp/astrolab%02d.json
	// syslog://localhost:5514/tag
	left, right, _ := strings.Cut(w.cfg.Out, ":")
	switch left {
	case "file":
		_, err = w.SenMl2File(right)
	case "syslog":
		connection := strings.TrimPrefix(right, "//")
		err = w.SenMl2Syslog(connection)
	default:
		err = ErrUnknownOut
	}
	if w.cfg.debug {
		slog.Debug("senMlWriter Write", "finished", "error", fmt.Sprintf("%v", err))
	}
	return err
}
