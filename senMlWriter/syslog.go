package senMlWriter

import (
	"log/syslog"
	"strings"

	"github.com/mainflux/senml"
)

const NoTag = "tag_is_not_set"

// SetPriority sets the priority (severity) of syslog.
// e.g: w.SetPriority(syslog.LOG_INFO|syslog.LOG_LOCAL7)
func (w *Writer) SetPriority(p syslog.Priority) *Writer {
	w.cfg.SyslogPriority = p
	return w
}

// SenMl2Syslog writes the senml.Pack to a syslog server
// with the given connection string. TCP supported only.
// It returns an error if the connection fails or the writing fails.
// The connection string should be in the form of "syslog://127.0.0.1:7814/senml02
func (w *Writer) SenMl2Syslog(connection string) error {

	tag := NoTag

	left, right, found := strings.Cut(connection, "/")
	if found && right != "" {
		connection = left
		tag = right
	}

	b, err := senml.Encode(w.p, senml.JSON)
	if err != nil {
		return err
	}

	return w.WriteToSyslog(connection, tag, b)
}

// WriteToSyslog writes the given data to a syslog server
// with the given connection string. TCP supported only.
// It returns an error if the connection fails or the writing fails.
// The connection string should be in the form of "syslog://127.0.0.1:7814/senml02".
func (w *Writer) WriteToSyslog(connection, tag string, data []byte) error {
	writer, err := syslog.Dial("tcp", connection, w.cfg.SyslogPriority, tag)

	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.Write(data)
	return err
}
