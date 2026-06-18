package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

type slowWriter struct {
	mu  sync.Mutex
	buf bytes.Buffer
	d   time.Duration
}

func (s *slowWriter) Write(p []byte) (int, error) {
	time.Sleep(s.d)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *slowWriter) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func TestFanoutWriter_MultiDestination(t *testing.T) {
	var a, b bytes.Buffer
	f := NewFanout(context.Background(), &a, &b)
	t.Cleanup(func() { _ = f.Close() })

	if _, err := f.Write([]byte("line\n")); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	if g := a.String(); g != "line\n" {
		t.Fatalf("a: got %q", g)
	}
	if g := b.String(); g != "line\n" {
		t.Fatalf("b: got %q", g)
	}
}

func TestFanoutWriter_WriteReturnsBeforeSlowBackend(t *testing.T) {
	slow := &slowWriter{d: 200 * time.Millisecond}
	cfg := DefaultFanoutConfig()
	cfg.WriteTimeout = time.Second
	f := NewFanoutConfig(context.Background(), cfg, slow)
	t.Cleanup(func() { _ = f.Close() })

	start := time.Now()
	if _, err := f.Write([]byte("x\n")); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("Write blocked %v", elapsed)
	}
	_ = f.Close()
	if slow.String() != "x\n" {
		t.Fatalf("got %q", slow.String())
	}
}

func TestFanoutWriter_DropWhenFull(t *testing.T) {
	slow := &slowWriter{d: 100 * time.Millisecond}
	cfg := DefaultFanoutConfig()
	cfg.Buffer = 1
	cfg.WriteTimeout = time.Second
	f := NewFanoutConfig(context.Background(), cfg, slow)
	t.Cleanup(func() { _ = f.Close() })

	_, _ = f.Write([]byte("first\n"))
	_, _ = f.Write([]byte("second\n"))
	_, _ = f.Write([]byte("third\n"))

	if f.Dropped() == 0 {
		t.Fatal("expected drops when buffer full")
	}
	_ = f.Close()
}

func TestFanoutWriter_BackendTimeoutBlacklist(t *testing.T) {
	slow := &slowWriter{d: 500 * time.Millisecond}
	var fast bytes.Buffer
	cfg := FanoutConfig{
		Buffer:       8,
		WriteTimeout: 20 * time.Millisecond,
		Cooldown:     200 * time.Millisecond,
	}
	f := NewFanoutConfig(context.Background(), cfg, slow, &fast)
	t.Cleanup(func() { _ = f.Close() })

	_, _ = f.Write([]byte("one\n"))
	deadline := time.Now().Add(time.Second)
	for f.BackendTimeouts() == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if f.BackendTimeouts() == 0 {
		t.Fatal("expected backend timeout")
	}
	if f.TimeoutDrops() == 0 {
		t.Fatal("expected timeout drop count")
	}

	_, _ = f.Write([]byte("two\n"))
	time.Sleep(30 * time.Millisecond)
	if !strings.Contains(fast.String(), "one\n") {
		t.Fatalf("fast backend: got %q", fast.String())
	}
	if strings.Contains(slow.String(), "two\n") {
		t.Fatalf("slow backend should be blacklisted, got %q", slow.String())
	}
}

func TestFanoutWriter_TimeoutShedsQueuedLine(t *testing.T) {
	slow := &slowWriter{d: 500 * time.Millisecond}
	cfg := FanoutConfig{
		Buffer:       4,
		WriteTimeout: 20 * time.Millisecond,
		Cooldown:     time.Second,
	}
	f := NewFanoutConfig(context.Background(), cfg, slow)
	t.Cleanup(func() { _ = f.Close() })

	_, _ = f.Write([]byte("first\n"))
	_, _ = f.Write([]byte("second\n"))

	deadline := time.Now().Add(time.Second)
	for f.BackendTimeouts() == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if f.Dropped() == 0 {
		t.Fatal("expected queued line shed after timeout")
	}
}

func TestFanoutWriter_WithLogConfig(t *testing.T) {
	var buf bytes.Buffer
	f := NewFanout(context.Background(), &buf)
	t.Cleanup(func() { _ = f.Close() })

	Config().Level(slog.LevelInfo).JSON().Out(f).Apply()
	Info("hello", "k", 1)
	_ = f.Close()

	out := buf.String()
	if !strings.Contains(out, `"msg":"hello"`) && !strings.Contains(out, `"msg": "hello"`) {
		t.Fatalf("got %q", out)
	}
}

func TestFanoutWriter_Closed(t *testing.T) {
	f := NewFanout(context.Background(), io.Discard)
	_ = f.Close()
	if _, err := f.Write([]byte("x")); err != io.ErrClosedPipe {
		t.Fatalf("got %v", err)
	}
}

func TestFanoutWriter_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer
	f := NewFanout(ctx, &buf)
	cancel()
	time.Sleep(20 * time.Millisecond)
	if _, err := f.Write([]byte("x\n")); err == nil {
		t.Fatal("expected error after ctx cancel")
	}
	_ = f.Close()
}
