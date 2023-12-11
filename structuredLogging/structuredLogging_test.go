package structuredLogging

import (
	"log"
	"log/slog"
	"os"
	"strings"
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
		"nats":   "nats://127.0.0.1:4222",
		"file":   "/tmp/do-log.log",
		"STDERR": true,
	})
	log.Println("Hallo World")
	slog.Error("this is my first error")
	slog.Debug("this is my first debug entry")
}

func TestGenerateLogfileName(t *testing.T) {
	fileName := GenerateLogfileName("/var/log/messenger-user.Current-2006-01.log")
	if strings.Contains(fileName, "user.Current") {
		t.Errorf("wrong filename %q", fileName)
	}
	if strings.Contains(fileName, "2006") {
		t.Errorf("wrong filename %q", fileName)
	}
}
