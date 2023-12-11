package structuredLogging

import (
	"log"
	"log/slog"
	"os"
	"testing"
)

func Test_StderrLogging(t *testing.T) {
	New(nil)
	child := slog.With(
		slog.Int("pid", os.Getpid()),
		slog.String("name", "test only"),
	)
	child.Error("child error")
	log.Println("Hallo World")
}

func Test_AllChannels(t *testing.T) {
	New(map[string]any{
		"nats":   "witest.itdesign.at:4222",
		"file":   "/var/log/do-log.log",
		"stderr": "yes",
	})
	log.Println("Hallo World")
	slog.Error("this is my first error")
	slog.Debug("this is my first debug entry")
}
