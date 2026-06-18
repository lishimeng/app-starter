package log

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestLevelFromString(t *testing.T) {
	cases := []struct {
		in   string
		want slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"WARNING", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"FINE", slog.LevelDebug},
	}
	for _, c := range cases {
		got, err := LevelFromString(c.in)
		if err != nil {
			t.Fatalf("%s: %v", c.in, err)
		}
		if got != c.want {
			t.Fatalf("%s: got %v want %v", c.in, got, c.want)
		}
	}
}

func TestFormatMessage_Printf(t *testing.T) {
	if got := formatMessage("n=%d", 1); got != "n=1" {
		t.Fatalf("got %q", got)
	}
}

func captureLogs(fn func()) string {
	var buf bytes.Buffer
	Config().Level(slog.LevelDebug).Out(&buf).Text().Apply()
	fn()
	return buf.String()
}

func TestInfo_Structured(t *testing.T) {
	out := captureLogs(func() {
		Info("query", "pageNum", 0, "pageSize", 10)
	})
	if !strings.Contains(out, "pageNum=0") || !strings.Contains(out, "pageSize=10") {
		t.Fatalf("expected structured fields, got: %q", out)
	}
}

func TestInfof_Printf(t *testing.T) {
	out := captureLogs(func() {
		Infof("web server listen %s", ":3000")
	})
	if !strings.Contains(out, "web server listen :3000") {
		t.Fatalf("got %q", out)
	}
}

func TestPkgInfo_AutoModule(t *testing.T) {
	out := captureLogs(func() {
		Info("hello")
	})
	if !strings.Contains(out, "module=log") {
		t.Fatalf("expected module=log, got %q", out)
	}
	if !strings.Contains(out, "hello") {
		t.Fatalf("expected hello, got %q", out)
	}
}

func TestFor_ExplicitModule(t *testing.T) {
	out := captureLogs(func() {
		For("syncdb").Info("created")
	})
	if !strings.Contains(out, "module=syncdb") {
		t.Fatalf("expected module=syncdb, got %q", out)
	}
}

func TestForEmpty_AutoModule(t *testing.T) {
	out := captureLogs(func() {
		For("").Info("x")
	})
	if !strings.Contains(out, "module=log") {
		t.Fatalf("expected module=log, got %q", out)
	}
}

func TestModuleFromFrame(t *testing.T) {
	frame := CallerFrame()
	mod := moduleFromFrame(frame)
	if mod != "log" {
		t.Fatalf("got %q want log", mod)
	}
}

func TestCallerFrame(t *testing.T) {
	frame := CallerFrame()
	if !strings.Contains(frame.File, "logger_test.go") {
		t.Fatalf("expected test file frame, got %q", frame.File)
	}
}

func TestSource_Disabled(t *testing.T) {
	var buf bytes.Buffer
	Config().Level(slog.LevelDebug).Out(&buf).Text().Caller(false).Apply()
	Info("hello")
	out := buf.String()
	if strings.Contains(out, "source=") {
		t.Fatalf("expected no source field, got: %q", out)
	}
}

func TestSource_Enabled(t *testing.T) {
	var buf bytes.Buffer
	Config().Level(slog.LevelDebug).Out(&buf).Text().Caller(true).Apply()
	Info("hello")
	out := buf.String()
	if !strings.Contains(out, "source=log/logger_test.go:") {
		t.Fatalf("expected source=file:line, got: %q", out)
	}
}

func TestSource_EnabledForExplicitModule(t *testing.T) {
	var buf bytes.Buffer
	Config().Level(slog.LevelDebug).Out(&buf).Text().Caller(true).Apply()
	l := For("syncdb")
	l.Info("created")
	out := buf.String()
	if !strings.Contains(out, "module=syncdb") {
		t.Fatalf("expected module=syncdb, got: %q", out)
	}
	if !strings.Contains(out, "source=log/logger_test.go:") {
		t.Fatalf("expected source=file:line, got: %q", out)
	}
}

func TestSetLevelFromString(t *testing.T) {
	if err := SetLevelFromString("ERROR"); err != nil {
		t.Fatal(err)
	}
	v, _ := globalLevel.Load().(*slog.LevelVar)
	if v.Level() != slog.LevelError {
		t.Fatalf("got %v", v.Level())
	}
	Config().Level(slog.LevelInfo).Text().Apply()
}

func TestWriteRaw(t *testing.T) {
	var buf bytes.Buffer
	Config().Out(&buf).Apply()
	WriteRaw("colored line")
	out := buf.String()
	if out != "colored line\n" {
		t.Fatalf("got %q", out)
	}
	if strings.Contains(out, "time=") {
		t.Fatal("WriteRaw should not use slog format")
	}
}

func TestApplySwitchTextToJSON(t *testing.T) {
	var buf bytes.Buffer
	Config().Out(&buf).Text().Apply()
	Info("text")
	Config().Out(&buf).JSON().Apply()
	Info("json")
	out := buf.String()
	if !strings.Contains(out, "time=") {
		t.Fatalf("expected text format line, got: %q", out)
	}
	if !strings.Contains(out, `"msg":"json"`) && !strings.Contains(out, `"msg": "json"`) {
		t.Fatalf("expected json format line, got: %q", out)
	}
}
