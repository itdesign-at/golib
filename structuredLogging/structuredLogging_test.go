package structuredLogging

import (
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
)

var funcWriteStdErr = reflect.ValueOf(writeStdErr).Pointer()
var funcWriteFile = reflect.ValueOf(writeFile).Pointer()
var funcWriteNats = reflect.ValueOf(writeNats).Pointer()

func Test_Params(t *testing.T) {

	var expectMap = map[string]map[string]uintptr{
		"xxwrong":                    {"STDERR": funcWriteStdErr},
		"":                           {"STDERR": funcWriteStdErr},
		"STDERR":                     {"STDERR": funcWriteStdErr},
		"/var/log/myLog.log":         {"/var/log/myLog.log": funcWriteFile},
		"nats://server.demo.at:4222": {"nats://server.demo.at:4222": funcWriteNats},
		"STDERR,/var/log/myLog.log,nats://server.demo.at:4222": {
			"STDERR":                     funcWriteStdErr,
			"/var/log/myLog.log":         funcWriteFile,
			"nats://server.demo.at:4222": funcWriteNats,
		},
	}

	sl := New()
	for param, expect := range expectMap {
		l := sl.Init(param)

		if len(sl.writers) != len(expect) {
			t.Errorf("test %q got wrong numbers of writers  %q (expected: %q)", param, len(sl.writers), len(expect))
		}

		for k, v := range expect {
			if f, ok := l.writers[k]; ok && reflect.ValueOf(f).Pointer() == v {
				t.Logf("OK param %q got writer %v", k, v)
			} else {
				t.Errorf("FAILED param %q doesn't get writer %v", param, v)
			}
		}

		for k, f := range l.writers {
			if _, ok := expect[k]; !ok {
				t.Errorf("FAILED param %q get wrong writer %v", param, reflect.ValueOf(f).Pointer())
			}
		}
	}

	param := []string{"STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222"}
	l := sl.Init(param...)
	if len(l.writers) != 3 {
		t.Errorf("FAILED param %q got wrong numbers of writers %q (expected: %q)", param, len(l.writers), 3)
	}

	l = sl.Init("STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222")
	if len(l.writers) != 3 {
		t.Errorf(
			"FAILED param \"STDERR\", "+
				"\"/var/log/myLog.log\", "+
				"\"nats://server.demo.at:4222\": "+
				"got wrong numbers of writers %q (expected: %q)", len(l.writers), 3)
	}

	l = sl.Init()
	if len(l.writers) == 1 && reflect.ValueOf(l.writers["STDERR"]).Pointer() == funcWriteStdErr {
		t.Logf("OK param %q got StdErr writer", "")
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", "")
	}

	var x []string
	l = sl.Init(x...)
	if len(l.writers) == 1 && reflect.ValueOf(l.writers["STDERR"]).Pointer() == funcWriteStdErr {
		t.Logf("OK param %q got StdErr writer", x)
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", x)
	}

	l = sl.Init(nil...)
	if len(l.writers) == 1 && reflect.ValueOf(l.writers["STDERR"]).Pointer() == funcWriteStdErr {
		t.Logf("OK param %q got StdErr writer", "nil")
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", "nil")
	}
}

func Test_StderrLogging(t *testing.T) {
	sl := New().Init()
	_ = sl //only for debugging

	child := slog.With(
		slog.Int("pid", os.Getpid()),
		slog.String("name", "test only"),
	)
	child.Error("child error")
	log.Println("Hallo World")
}

func Test_AllChannels(t *testing.T) {
	//	sl := New().Init("nats://127.0.0.1:4222", "/tmp/do-log.log", "STDERR")
	sl := New().Init("nats://127.0.0.1:4222")
	_ = sl //only for debugging

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
