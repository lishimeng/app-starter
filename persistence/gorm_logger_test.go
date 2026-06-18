package persistence

import (
	"bytes"
	"strings"
	"testing"
	"time"

	applog "github.com/lishimeng/app-starter/log"
	gormlogger "gorm.io/gorm/logger"
)

func TestSlogGormLoggerTrace(t *testing.T) {
	var buf bytes.Buffer
	applog.Config().Output(&buf).LevelFromString("debug").Apply()

	l := newSlogGormLogger(gormlogger.Info)
	l.Trace(nil, time.Now().Add(-5*time.Millisecond), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	out := buf.String()
	if strings.Contains(out, "time=") {
		t.Fatalf("gorm sql should bypass slog text format, got: %q", out)
	}
	if !strings.Contains(out, "SELECT 1") {
		t.Fatalf("expected sql text, got: %q", out)
	}
	if !strings.Contains(out, "[rows:1]") {
		t.Fatalf("expected rows in gorm format, got: %q", out)
	}
}

func TestSlogGormLoggerSilent(t *testing.T) {
	var buf bytes.Buffer
	applog.Config().Output(&buf).LevelFromString("debug").Apply()

	l := newSlogGormLogger(gormlogger.Silent)
	l.Trace(nil, time.Now(), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	if buf.Len() != 0 {
		t.Fatalf("expected no output in silent mode, got: %q", buf.String())
	}
}

func TestGormLogLevel(t *testing.T) {
	if gormLogLevel(false) != gormlogger.Silent {
		t.Fatal("expected silent when debug disabled")
	}
	if gormLogLevel(true) != gormlogger.Info {
		t.Fatal("expected info when debug enabled")
	}
}
