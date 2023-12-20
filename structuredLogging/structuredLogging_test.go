package structuredLogging

import (
	"github.com/itdesign-at/golib/keyvalue"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_Params(t *testing.T) {

	testWriter := func(param string, s *SlogLogger, e map[string]uintptr) {
		if len(s.writers) != len(e) {
			t.Errorf("test %q got wrong numbers of writers  %q (expected: %q)", param, len(s.writers), len(e))
		}

		for k, v := range e {
			if f, ok := s.writers[k]; ok && reflect.ValueOf(f).Pointer() == v {
				t.Logf("OK param %q got writer %v", k, v)
			} else {
				t.Errorf("FAILED param %q doesn't get writer %v", param, v)
			}
		}

		for k, f := range s.writers {
			if _, ok := e[k]; !ok {
				t.Errorf("FAILED param %q get wrong writer %v", param, reflect.ValueOf(f).Pointer())
			}
		}
	}

	param := "xxwrong"
	sl := New(param).Init()
	expect := map[string]uintptr{"STDERR": reflect.ValueOf(sl.writeStdErr).Pointer()}
	testWriter(param, sl, expect)

	param = ""
	sl = New(param).Init()
	expect = map[string]uintptr{"STDERR": reflect.ValueOf(sl.writeStdErr).Pointer()}
	testWriter(param, sl, expect)

	param = "STDERR"
	sl = New(param).Init()
	expect = map[string]uintptr{"STDERR": reflect.ValueOf(sl.writeStdErr).Pointer()}
	testWriter(param, sl, expect)

	param = "/var/log/myLog.log"
	sl = New(param).Init()
	expect = map[string]uintptr{"/var/log/myLog.log": reflect.ValueOf(sl.writeFile).Pointer()}
	testWriter(param, sl, expect)

	param = "nats://server.demo.at:4222"
	sl = New(param).Init()
	expect = map[string]uintptr{"nats://server.demo.at:4222": reflect.ValueOf(sl.writeNats).Pointer()}
	testWriter(param, sl, expect)

	param = "STDERR,/var/log/myLog.log,nats://server.demo.at:4222"
	sl = New(param).Init()
	expect = map[string]uintptr{
		"STDERR":                     reflect.ValueOf(sl.writeStdErr).Pointer(),
		"/var/log/myLog.log":         reflect.ValueOf(sl.writeFile).Pointer(),
		"nats://server.demo.at:4222": reflect.ValueOf(sl.writeNats).Pointer()}
	testWriter(param, sl, expect)

	paramSlice := []string{"STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222"}
	sl = New(paramSlice...).Init()
	if len(sl.writers) != 3 {
		t.Errorf("FAILED param %q got wrong numbers of writers %q (expected: %q)", param, len(sl.writers), 3)
	}

	sl = New("STDERR", "/var/log/myLog.log", "nats://server.demo.at:4222").Init()
	if len(sl.writers) != 3 {
		t.Errorf(
			"FAILED param \"STDERR\", "+
				"\"/var/log/myLog.log\", "+
				"\"nats://server.demo.at:4222\": "+
				"got wrong numbers of writers %q (expected: %q)", len(sl.writers), 3)
	}

	sl = New().Init()
	if len(sl.writers) == 1 && reflect.ValueOf(sl.writers["STDERR"]).Pointer() == reflect.ValueOf(sl.writeStdErr).Pointer() {
		t.Logf("OK param %q got StdErr writer", "")
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", "")
	}

	var x []string
	sl = New(x...).Init()
	if len(sl.writers) == 1 && reflect.ValueOf(sl.writers["STDERR"]).Pointer() == reflect.ValueOf(sl.writeStdErr).Pointer() {
		t.Logf("OK param %q got StdErr writer", x)
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", x)
	}

	sl = New(nil...).Init()
	if len(sl.writers) == 1 && reflect.ValueOf(sl.writers["STDERR"]).Pointer() == reflect.ValueOf(sl.writeStdErr).Pointer() {
		t.Logf("OK param %q got StdErr writer", "nil")
	} else {
		t.Errorf("FAILED  param %q doesn't get stdErr writer", "nil")
	}

	sl = New().Init()
	if len(sl.writers) == 1 && reflect.ValueOf(sl.writers["STDERR"]).Pointer() == reflect.ValueOf(sl.writeStdErr).Pointer() {
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
	sl := New("nats://127.0.0.1:4222", "/tmp/do-log.log", "STDERR").Init()
	//sl := New("nats://127.0.0.1:4222").Init()
	_ = sl //only for debugging

	log.Println("Hallo World")
	slog.Error("this is my first error")
	slog.Debug("this is my first debug entry")
}

func Test_Levels(t *testing.T) {
	sl := New("STDERR").Parameter(keyvalue.Record{"Level": "error"}).Init()
	_ = sl //only for debugging

	log.Println("log.Println() is a info entry")
	slog.Info("this is my first info entry")
	slog.Warn("this is my first warning entry")
	slog.Error("this is my first error entry")
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

func Test_StdOut(t *testing.T) {
	sl := New("StdOut").Init()
	_ = sl //only for debugging

	log.Println("log.Println() is a info entry")
	slog.Info("this is my first info entry")
	slog.Warn("this is my first warning entry")
	slog.Error("this is my first error entry")
	slog.Debug("this is my first debug entry")
}
