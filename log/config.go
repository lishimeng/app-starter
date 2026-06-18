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
	defaultHandler atomic.Value // handlerSlot
	globalLevel    atomic.Value // *slog.LevelVar
	showSource     atomic.Bool
	outMu          sync.RWMutex
	defaultOut     io.Writer = os.Stdout
	outWriter      io.Writer = os.Stdout
	rawMu          sync.Mutex
)

type handlerSlot struct {
	h slog.Handler
}

func storeHandler(h slog.Handler) {
	defaultHandler.Store(handlerSlot{h: h})
}

func loadHandler() slog.Handler {
	if slot, ok := defaultHandler.Load().(handlerSlot); ok && slot.h != nil {
		return slot.h
	}
	return slog.Default().Handler()
}

func init() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	globalLevel.Store(lvl)
	outWriter = defaultOut
	storeHandler(newHandler(defaultOut, lvl, false))
}

// Options builds slog handler settings (chainable).
type Options struct {
	level      *slog.LevelVar
	writer     io.Writer
	json       bool
	withCaller bool
}

// Config starts a configuration chain.
func Config() *Options {
	lvl, _ := globalLevel.Load().(*slog.LevelVar)
	if lvl == nil {
		lvl = new(slog.LevelVar)
	}
	return &Options{
		level:  lvl,
		writer: defaultOut,
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

func (c *Options) Out(w io.Writer) *Options {
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

// Caller enables source file:line on each record; entries attach source= via stack walk (default off).
func (c *Options) Caller(enabled bool) *Options {
	if c != nil {
		c.withCaller = enabled
	}
	return c
}

// Apply installs handler as slog.Default and package default.
func (c *Options) Apply() {
	if c == nil {
		return
	}
	globalLevel.Store(c.level)
	showSource.Store(c.withCaller)
	storeHandler(newHandler(c.writer, c.level, c.json))
	if c.writer != nil {
		outMu.Lock()
		outWriter = c.writer
		outMu.Unlock()
	}
	slog.SetDefault(slog.New(currentHandler()))
}

// WriteRaw writes pre-formatted text directly to the configured output (no slog key=value wrapper).
// Use for ANSI-colored lines such as GORM SQL traces.
func WriteRaw(msg string) {
	outMu.RLock()
	w := outWriter
	outMu.RUnlock()
	if w == nil {
		w = defaultOut
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
	return loadHandler()
}

func sourceEnabled() bool {
	return showSource.Load()
}
