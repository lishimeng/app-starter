package log

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultFanoutBuffer   = 256
	defaultWriteTimeout   = 2 * time.Second
	defaultBlacklistAfter = 30 * time.Second
)

// FanoutConfig controls async fan-out behavior.
type FanoutConfig struct {
	Buffer       int           // channel capacity; default 256
	WriteTimeout time.Duration // per-backend write deadline; default 2s
	Cooldown     time.Duration // blacklist duration after timeout; default 30s
}

func DefaultFanoutConfig() FanoutConfig {
	return FanoutConfig{
		Buffer:       defaultFanoutBuffer,
		WriteTimeout: defaultWriteTimeout,
		Cooldown:     defaultBlacklistAfter,
	}
}

type fanoutBackend struct {
	w              io.Writer
	blacklistUntil atomic.Int64 // unix nano; 0 = active
}

func (b *fanoutBackend) blacklisted(now time.Time) bool {
	until := b.blacklistUntil.Load()
	if until == 0 {
		return false
	}
	return now.Before(time.Unix(0, until))
}

func (b *fanoutBackend) blacklist(now time.Time, cooldown time.Duration) {
	if cooldown <= 0 {
		return
	}
	b.blacklistUntil.Store(now.Add(cooldown).UnixNano())
}

// FanoutWriter duplicates log lines to multiple io.Writer backends asynchronously.
type FanoutWriter struct {
	ctx    context.Context
	cancel context.CancelFunc

	ch       chan []byte
	backends []*fanoutBackend
	cfg      FanoutConfig

	wg     sync.WaitGroup
	close  sync.Once
	closed atomic.Bool

	dropped         atomic.Uint64
	timeoutDrops    atomic.Uint64
	backendTimeouts atomic.Uint64
}

// NewFanout returns a fan-out writer with DefaultFanoutConfig.
func NewFanout(ctx context.Context, writers ...io.Writer) *FanoutWriter {
	return NewFanoutConfig(ctx, DefaultFanoutConfig(), writers...)
}

// NewFanoutBuffer is deprecated naming; use NewFanoutConfig.
func NewFanoutBuffer(size int, writers ...io.Writer) *FanoutWriter {
	cfg := DefaultFanoutConfig()
	cfg.Buffer = size
	return NewFanoutConfig(context.Background(), cfg, writers...)
}

// NewFanoutConfig builds a fan-out writer. ctx cancellation stops the worker.
func NewFanoutConfig(ctx context.Context, cfg FanoutConfig, writers ...io.Writer) *FanoutWriter {
	if ctx == nil {
		ctx = context.Background()
	}
	if cfg.Buffer <= 0 {
		cfg.Buffer = defaultFanoutBuffer
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = defaultWriteTimeout
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = defaultBlacklistAfter
	}

	ctx, cancel := context.WithCancel(ctx)
	backends := make([]*fanoutBackend, 0, len(writers))
	for _, w := range writers {
		if w != nil {
			backends = append(backends, &fanoutBackend{w: w})
		}
	}
	f := &FanoutWriter{
		ctx:      ctx,
		cancel:   cancel,
		ch:       make(chan []byte, cfg.Buffer),
		backends: backends,
		cfg:      cfg,
	}
	f.wg.Add(1)
	go f.run()
	return f
}

func (f *FanoutWriter) run() {
	defer f.wg.Done()
	for {
		select {
		case <-f.ctx.Done():
			return
		case line, ok := <-f.ch:
			if !ok {
				return
			}
			f.dispatch(line)
		}
	}
}

func (f *FanoutWriter) dispatch(line []byte) {
	now := time.Now()
	for _, b := range f.backends {
		if b.blacklisted(now) {
			continue
		}
		if err := writeWithTimeout(f.ctx, b.w, line, f.cfg.WriteTimeout); err != nil {
			f.backendTimeouts.Add(1)
			f.timeoutDrops.Add(1)
			b.blacklist(now, f.cfg.Cooldown)
			f.dropOneQueued()
		}
	}
}

func writeWithTimeout(ctx context.Context, w io.Writer, p []byte, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := w.Write(p)
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (f *FanoutWriter) dropOneQueued() {
	select {
	case <-f.ch:
		f.dropped.Add(1)
	default:
	}
}

// Write enqueues a log line copy; returns immediately without waiting for backends.
func (f *FanoutWriter) Write(p []byte) (int, error) {
	if f == nil {
		return 0, io.ErrClosedPipe
	}
	if f.closed.Load() {
		return 0, io.ErrClosedPipe
	}
	if err := f.ctx.Err(); err != nil {
		return 0, err
	}
	if len(p) == 0 {
		return 0, nil
	}
	if len(f.backends) == 0 {
		return len(p), nil
	}
	buf := append([]byte(nil), p...)
	select {
	case f.ch <- buf:
	case <-f.ctx.Done():
		return 0, f.ctx.Err()
	default:
		f.dropped.Add(1)
	}
	return len(p), nil
}

// Dropped returns lines dropped (channel full or shed after backend timeout).
func (f *FanoutWriter) Dropped() uint64 {
	if f == nil {
		return 0
	}
	return f.dropped.Load()
}

// TimeoutDrops returns lines not delivered due to backend write timeouts.
func (f *FanoutWriter) TimeoutDrops() uint64 {
	if f == nil {
		return 0
	}
	return f.timeoutDrops.Load()
}

// BackendTimeouts returns how many backend writes hit the timeout.
func (f *FanoutWriter) BackendTimeouts() uint64 {
	if f == nil {
		return 0
	}
	return f.backendTimeouts.Load()
}

// Close drains pending lines and stops the worker.
func (f *FanoutWriter) Close() error {
	if f == nil {
		return nil
	}
	f.close.Do(func() {
		f.closed.Store(true)
		close(f.ch)
		f.wg.Wait()
		f.cancel()
	})
	return nil
}
