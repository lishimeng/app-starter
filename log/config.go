package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
)

var (
	defaultHandler atomic.Value // slog.Handler
	globalLevel    atomic.Value // slog.LevelVar
	outMu          sync.RWMutex
	outWriter      io.Writer = os.Stderr
	rawMu          sync.Mutex
)

func init() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	globalLevel.Store(lvl)
	outWriter = os.Stderr
	defaultHandler.Store(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}

// Options builds slog handler settings (chainable).
type Options struct {
	level  *slog.LevelVar
	writer io.Writer
	json   bool
}

// Config starts a configuration chain.
func Config() *Options {
	lvl, _ := globalLevel.Load().(*slog.LevelVar)
	if lvl == nil {
		lvl = new(slog.LevelVar)
	}
	return &Options{
		level:  lvl,
		writer: os.Stderr,
	}
}

func (c *Options) Level(l slog.Level) *Options {
	if c != nil && c.level != nil {
		c.level.Set(l)
	}
	return c
}

func (c *Options) LevelFromString(s string) *Options {
	if c == nil {
		return c
	}
	lvl, err := LevelFromString(s)
	if err == nil {
		c.level.Set(lvl)
	}
	return c
}

func (c *Options) Output(w io.Writer) *Options {
	if c != nil && w != nil {
		c.writer = w
	}
	return c
}

func (c *Options) Text() *Options {
	if c != nil {
		c.json = false
	}
	return c
}

func (c *Options) JSON() *Options {
	if c != nil {
		c.json = true
	}
	return c
}

// Apply installs handler as slog.Default and package default.
func (c *Options) Apply() {
	if c == nil {
		return
	}
	opts := &slog.HandlerOptions{Level: c.level}
	var h slog.Handler
	if c.json {
		h = slog.NewJSONHandler(c.writer, opts)
	} else {
		h = slog.NewTextHandler(c.writer, opts)
	}
	globalLevel.Store(c.level)
	defaultHandler.Store(h)
	if c.writer != nil {
		outMu.Lock()
		outWriter = c.writer
		outMu.Unlock()
	}
	slog.SetDefault(slog.New(h))
}

// WriteRaw writes pre-formatted text directly to the configured output (no slog key=value wrapper).
// Use for ANSI-colored lines such as GORM SQL traces.
func WriteRaw(msg string) {
	outMu.RLock()
	w := outWriter
	outMu.RUnlock()
	if w == nil {
		w = os.Stderr
	}
	rawMu.Lock()
	defer rawMu.Unlock()
	fmt.Fprintln(w, msg)
}

// SetLevelFromString changes global log level at runtime.
func SetLevelFromString(s string) error {
	lvl, err := LevelFromString(s)
	if err != nil {
		return err
	}
	if v, ok := globalLevel.Load().(*slog.LevelVar); ok && v != nil {
		v.Set(lvl)
	}
	return nil
}

func currentHandler() slog.Handler {
	if h, ok := defaultHandler.Load().(slog.Handler); ok && h != nil {
		return h
	}
	return slog.Default().Handler()
}
