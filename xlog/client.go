package xlog

import (
	"fmt"
	"github.com/itdesign-at/golib/converter"
	"github.com/nats-io/nats.go"
	"io"
	"log/slog"
	"math"
	"os"
	"strings"
	"sync"
)

const (
	LevelAudit  = slog.Level(2)
	KeyCategory = "category"
	KeyAction   = "action"
	KeyUser     = "user"
	KeyError    = "error"
	KeyExtended = "extended"
)

var LevelNames = map[slog.Leveler]string{
	LevelAudit: "AUDIT",
}

// Category returns a slog attribute with the given category.
func Category(category string) slog.Attr {
	return slog.String(KeyCategory, category)
}

// Action returns a slog attribute with the given action.
func Action(action string) slog.Attr {
	return slog.String(KeyAction, action)
}

// User returns a slog attribute with the given user.
func User(user string) slog.Attr {
	return slog.String(KeyUser, user)
}

// Error returns a slog attribute with the given error.
func Error(err error) slog.Attr {
	return slog.Any(KeyError, err)
}

// Extended returns a slog attribute with the given arguments.
func Extended(args ...any) slog.Attr {
	if len(args) == 1 {
		return slog.Any(KeyExtended, args[0])
	} else {
		return slog.Group(KeyExtended, args...)
	}
}

// ClientOptions holds the configuration options for the mlog logging client.
type ClientOptions struct {

	// log destinations
	// supported: stdout, stderr, /path/to/logfile, nats://host:port
	Destinations []string

	// Log Level (default: audit)
	Level string

	// AddSource add go source file name and line number to log output
	AddSource bool

	// Hostname added to log output and NATS subject (default: os.Hostname())
	Hostname string

	// optional override Stdout and Stderr writers for testing (default: os.Stdout, os.Stderr)
	Stdout io.Writer
	Stderr io.Writer
}

// parseLevel parses the log level from the options Level string.
// If the string is not a valid log level, it returns the info log level.
func (o *ClientOptions) parseLevel() slog.Level {

	// check if level is a custom level name
	for level, name := range LevelNames {
		if strings.EqualFold(o.Level, name) {
			return level.Level()
		}
	}

	// use slog.Level default parser
	var l slog.Level
	var err = l.UnmarshalText([]byte(o.Level))
	if err != nil {
		l = LevelAudit
	}
	return l
}

var (
	defaultHostname, _ = os.Hostname()
	defaultOptions     = ClientOptions{
		Level:    LevelNames[LevelAudit],
		Hostname: defaultHostname,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
	}
)

// Client is the heimdall logging client that writes to multiple destinations.
// It supports writing to stdout, stderr, logfiles and nats.
// The client is initialized with a slog logger that writes to all configured destinations.
type Client struct {

	// client options
	options ClientOptions

	// writer is an io.Writer that writes to all configured destinations
	writer *multiWriter

	// the initialized slog logger
	logger *slog.Logger

	// open logfiles
	logfiles []io.WriteCloser

	// open nats connections
	connections []*nats.Conn
}

// NewClient creates a new heimdall logging client with the given options.
// If options are not set, default options are used.
func NewClient(options ClientOptions) (*Client, error) {

	// set default options
	if options.Level == "" {
		options.Level = defaultOptions.Level
	}

	if options.Hostname == "" {
		options.Hostname = defaultOptions.Hostname
	}

	if options.Stdout == nil {
		options.Stdout = defaultOptions.Stdout
	}

	if options.Stderr == nil {
		options.Stderr = defaultOptions.Stderr
	}

	client := &Client{
		options:     options,
		writer:      newMultiWriter(),
		logfiles:    make([]io.WriteCloser, 0),
		connections: make([]*nats.Conn, 0),
	}

	// initialize client
	err := client.init()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// init initializes the heimdall logging client with the configured options.
// It sets up the log writers, the slog logger and the Nats connection.
func (c *Client) init() error {

	for _, dest := range c.options.Destinations {
		switch {
		// stdout
		case strings.ToLower(dest) == "stdout":
			c.writer.append(c.options.Stdout)
		// stderr
		case strings.ToLower(dest) == "stderr":
			c.writer.append(c.options.Stderr)
		// nats
		case strings.HasPrefix(strings.ToLower(dest), "nats://"):
			// connect to nats server
			connection, err := nats.Connect(dest,
				nats.RetryOnFailedConnect(true),
				nats.MaxReconnects(math.MaxInt))
			if err != nil {
				return err
			}
			c.connections = append(c.connections, connection)

			// create nats writer
			nw := newNatsWriter(connection, fmt.Sprintf("log.%s.%s", c.options.Level, converter.Normalize(c.options.Hostname)))
			c.writer.append(nw)
		// logfile
		default:
			logfile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			c.writer.append(logfile)
			c.logfiles = append(c.logfiles, logfile)
		}

	}

	// create log format handler
	h := slog.NewJSONHandler(
		c.writer,
		&slog.HandlerOptions{
			AddSource: c.options.AddSource,
			Level:     c.options.parseLevel(),
			// replace log level with custom level names
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level, ok := a.Value.Any().(slog.Level)
					if ok {
						levelLabel, exists := LevelNames[level]
						if !exists {
							levelLabel = level.String()
						}
						a.Value = slog.StringValue(levelLabel)
					}
				}
				return a
			},
		},
	)

	// create logger
	c.logger = slog.New(h)

	// always add node to logger
	c.logger = c.logger.With("hostname", c.options.Hostname)

	// set slog logger as default logger
	slog.SetDefault(c.logger)

	return nil
}

// Close closes all open logfiles and the Nats connection.
// It should be called in the main function before the program exits.
func (c *Client) Close() error {

	// close open logfiles
	for _, logfile := range c.logfiles {
		err := logfile.Close()
		if err != nil {
			return err
		}
	}

	// close nats connection
	for _, connection := range c.connections {
		err := connection.Flush()
		if err != nil {
			return err
		}
		connection.Close()
	}

	return nil
}

// natsWriter is a writer that writes to a Nats connection.
// It is used to write log messages to a Nats server.
type natsWriter struct {

	// the Nats subject to write to
	subject string

	// the Nats connection
	connection *nats.Conn
}

func newNatsWriter(connection *nats.Conn, subject string) *natsWriter {
	return &natsWriter{connection: connection, subject: subject}
}

// Write writes the given bytes to the Nats connection.
// It returns the number of bytes written and an error if the write fails.
func (w *natsWriter) Write(p []byte) (int, error) {
	err := w.connection.Publish(w.subject, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// multiWriter is a writer that writes to multiple io.Writers.
// It is used to write log messages to multiple destinations.
type multiWriter struct {

	// writers holds all writers
	writers []io.Writer
}

// newMultiWriter creates a new multiWriter with the given writers.
func newMultiWriter() *multiWriter {
	return &multiWriter{writers: make([]io.Writer, 0)}
}

// append adds the given writer to the multiWriter.
func (mw *multiWriter) append(w io.Writer) {
	mw.writers = append(mw.writers, w)
}

// This is the implementation of the io.Writer interface.
// It writes the given bytes to all configured io.Writers.
// It waits for all writers to finish before returning.
// All errors are ignored.
func (mw *multiWriter) Write(p []byte) (int, error) {
	var wg sync.WaitGroup

	for _, w := range mw.writers {
		wg.Add(1)
		go func(w io.Writer, p []byte) {
			_, _ = w.Write(p)
			wg.Done()
		}(w, p)
	}

	wg.Wait()
	return len(p), nil
}
