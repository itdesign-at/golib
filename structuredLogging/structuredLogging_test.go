package structuredLogging

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/itdesign-at/golib/keyvalue"
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

	logFilename := filepath.Join("/tmp", GenerateLogfileName("test-messenger-user.Current-2006-01.log"))
	New(logFilename).Init()
	slog.Info("just a test")
	entries, _ := os.ReadDir("/tmp")
	var found bool
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "test-messenger-") && strings.HasSuffix(e.Name(), ".log") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("generated file %q not found", logFilename)
	}
	content, _ := os.ReadFile(logFilename)
	if !bytes.Contains(content, []byte("just a test")) {
		t.Errorf("content of file %q is wrong", logFilename)
	}
	_ = os.Remove(logFilename)

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

func Test_Nats(t *testing.T) {

	var t0 time.Time

	// Werner: ausnahmsweise mit fmt.Println die Zeiten raus schreiben,
	// das das t.Log beim Testen umgeaendert wird.
	t0 = time.Now()
	sl := New("nats://127.0.0.1").Init()
	fmt.Println("NATS init: ", time.Now().Sub(t0).String())

	_ = sl //only for debugging

	t0 = time.Now()
	log.Println("log.Println() is a info entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Info("this is my first info entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Warn("this is my first warning entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Error("this is my first error entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Debug("this is my first debug entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	withSubject := New("nats://127.0.0.1").Parameter(
		map[string]interface{}{"natsSubject": "messages.watchit"}).Init()
	fmt.Println("NATS init: ", time.Now().Sub(t0).String())

	_ = withSubject //only for debugging

	t0 = time.Now()
	log.Println("log.Println() is a info entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Info("this is my first info entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Warn("this is my first warning entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Error("this is my first error entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())

	t0 = time.Now()
	slog.Debug("this is my first debug entry")
	fmt.Println("NATS logging: ", time.Now().Sub(t0).String())
}

func Test_PrepareNatsSubject(t *testing.T) {
	sl := New("nats://127.0.0.1").Init()
	var expected = map[string]string{
		"":                   "slog.UNKNOWN",
		"x":                  "slog.UNKNOWN",
		"slog.UNKNOWN":       "slog.UNKNOWN",
		`level ERROR`:        "slog.UNKNOWN",
		`{"level": "ERROR"}`: "slog.ERROR", // valid JSON must work
	}
	for k, v := range expected {
		subj := sl.prepareNatsSubject([]byte(k))
		if subj != v {
			t.Errorf("wrong subject - expected %q but got %q", v, subj)
		}
	}

	var logMessage = []byte(`{"msg": "just to test","level": "DEBUG"}`)

	expected = map[string]string{
		"":                   "slog.DEBUG", // OK -> default subject
		"mySubject.test":     "mySubject.test",
		`{"level": "ERROR"}`: `{"level": "ERROR"}`,
		"mySubject.LOGLEVEL": "mySubject.DEBUG", // OK -> custom subject
		"LOGLEVEL.LOGLEVEL":  "DEBUG.DEBUG",     // OK -> custom subject
	}

	for sbjTemplate, v := range expected {
		sl = New("nats://127.0.0.1").Parameter(map[string]interface{}{
			"NatsSubject": sbjTemplate,
		}).Init()
		subj := sl.prepareNatsSubject(logMessage)
		if subj != v {
			t.Errorf("wrong subject - entered %q expected %q but got %q", sbjTemplate, v, subj)
		}
	}

	sl = New("nats://127.0.0.1/mySubject.LOGLEVEL").Init()
	subj := sl.prepareNatsSubject(logMessage)
	if subj != "mySubject.DEBUG" {
		t.Errorf("wrong subject expected mySubject.DEBUG derived from URL param, got %q", subj)
	}
}
