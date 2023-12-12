package structuredLogging

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func Test_Params(t *testing.T) {

	var expectMap = map[string]string{
		"xxwrong":                    "STDERR;yes",
		"":                           "STDERR;yes",
		"STDERR":                     "STDERR;yes",
		"/var/log/myLog.log":         "file;/var/log/myLog.log",
		"nats://server.demo.at:4222": "nats;nats://server.demo.at:4222",
	}

	for k, v := range expectMap {
		l := New(k)
		left, right, _ := strings.Cut(v, ";")
		if len(l.myConfig) == 1 && l.myConfig[left] == right {
			t.Logf("OK sent %q and got %q:%q", k, left, right)
		} else {
			t.Errorf("wrong param sent %q and got %q:%q", k, left, right)
		}
		asSlice := []string{k}
		l = New(asSlice)
		if len(l.myConfig) == 1 && l.myConfig[left] == right {
			t.Logf("OK sent []string{%q} and got %q:%q", k, left, right)
		} else {
			t.Errorf("wrong param sent %q and got %q:%q", k, left, right)
		}
	}

	l := New("STDERR,/var/log/myLog.log,nats://server.demo.at:4222")
	if len(l.myConfig) != 3 {
		t.Errorf("wrong params")
	}

	l = New([]string{"STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222"})
	if len(l.myConfig) != 3 {
		t.Errorf("wrong params")
	}

}

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
