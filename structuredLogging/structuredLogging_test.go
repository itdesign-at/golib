package structuredLogging

import (
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
)

//TODO: die Tests sind noch mist! muss ich noch fertig machen!

var funcWriteStdErr = reflect.ValueOf(writeStdErr).Pointer()
var funcWriteFile = reflect.ValueOf(writeFile).Pointer()
var funcWriteNats = reflect.ValueOf(writeNats).Pointer()

func Test_Params(t *testing.T) {

	var expectMap = map[string][]uintptr{
		"xxwrong":                    {funcWriteStdErr},
		"":                           {funcWriteStdErr},
		"STDERR":                     {funcWriteStdErr},
		"/var/log/myLog.log":         {funcWriteFile},
		"nats://server.demo.at:4222": {funcWriteNats},
		"STDERR,/var/log/myLog.log,nats://server.demo.at:4222": {funcWriteStdErr, funcWriteFile, funcWriteNats},
	}

	for k, v := range expectMap {
		l := New(k)
		if f := l.writers[k]; len(l.writers) == 1 && reflect.ValueOf(f).Pointer() == v[0] {
			t.Logf("OK sent %q and got %q:%q", k, v[0], reflect.ValueOf(l.writers[k]).Pointer())
		} else {
			t.Errorf("wrong param sent %q and got %q:%q", k, v[0], reflect.ValueOf(l.writers[k]).Pointer())
		}
	}

	l := New([]string{"STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222"}...)
	if len(l.writers) != 3 {
		t.Errorf("wrong params")
	}

	l = New("STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222")
	if len(l.writers) != 3 {
		t.Errorf("wrong params")
	}

	l = New()
	if len(l.writers) == 1 && reflect.ValueOf(l.writers["STDERR"]).Pointer() == funcWriteStdErr {
		t.Logf("OK sent %q and got %q:%q", "nil", funcWriteStdErr, reflect.ValueOf(l.writers["STDERR"]).Pointer())
	} else {
		t.Errorf("wrong param sent %q and got %q:%q", "nil", funcWriteStdErr, reflect.ValueOf(l.writers["STDERR"]).Pointer())
	}

	var x []string
	l = New(x...)
	if len(l.writers) == 1 && reflect.ValueOf(l.writers["STDERR"]).Pointer() == funcWriteStdErr {
		t.Logf("OK sent %q and got %q:%q", "nil", funcWriteStdErr, reflect.ValueOf(l.writers["STDERR"]).Pointer())
	} else {
		t.Errorf("wrong param sent %q and got %q:%q", "nil", funcWriteStdErr, reflect.ValueOf(l.writers["STDERR"]).Pointer())
	}

}

func Test_StderrLogging(t *testing.T) {

	New()
	child := slog.With(
		slog.Int("pid", os.Getpid()),
		slog.String("name", "test only"),
	)
	child.Error("child error")
	log.Println("Hallo World")
}

func Test_AllChannels(t *testing.T) {
	l := New("nats://127.0.0.1:4222", "/tmp/do-log.log", "STDERR")
	_ = l
	log.Println("Hallo World")
	slog.Error("this is my first error")
	slog.Debug("this is my first debug entry")
}

/*
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

*/

func TestGenerateLogfileName(t *testing.T) {
	fileName := GenerateLogfileName("/var/log/messenger-user.Current-2006-01.log")
	if strings.Contains(fileName, "user.Current") {
		t.Errorf("wrong filename %q", fileName)
	}
	if strings.Contains(fileName, "2006") {
		t.Errorf("wrong filename %q", fileName)
	}
}
