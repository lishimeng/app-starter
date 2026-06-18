package log

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// Logger wraps slog with optional fixed module and chainable With.
type Logger struct {
	inner  *slog.Logger
	module string
	attrs  []any
}

// For returns a logger with explicit module name. Empty module uses caller inference.
func For(module string) *Logger {
	return &Logger{
		inner:  slog.New(currentHandler()),
		module: module,
	}
}

func rootLogger() *Logger {
	return &Logger{inner: slog.New(currentHandler())}
}

// Slog returns a stdlib slog.Logger sharing the configured handler and fixed module.
func (l *Logger) Slog() *slog.Logger {
	base := slog.New(currentHandler())
	if l != nil && l.module != "" {
		return base.With(slog.String("module", l.module))
	}
	return base
}

// Slog returns a stdlib slog.Logger for integrations (e.g. GORM).
func Slog(module string) *slog.Logger {
	return For(module).Slog()
}

// With returns a logger with extra attributes (module still auto or fixed).
func With(args ...any) *Logger {
	return rootLogger().With(args...)
}

func (l *Logger) With(args ...any) *Logger {
	if l == nil {
		return rootLogger().With(args...)
	}
	next := &Logger{
		inner:  l.inner,
		module: l.module,
		attrs:  append(append([]any{}, l.attrs...), args...),
	}
	return next
}

func (l *Logger) prependModule(frame runtime.Frame, attrs []any) []any {
	if l == nil {
		l = rootLogger()
	}
	if l.module != "" {
		return append([]any{slog.String("module", l.module)}, attrs...)
	}
	if mod := moduleFromFrame(frame); mod != "" {
		return append([]any{slog.String("module", mod)}, attrs...)
	}
	return attrs
}

func (l *Logger) handle(level slog.Level, msg string, attrs []any) *Logger {
	if l == nil {
		l = rootLogger()
	}
	ctx := context.Background()
	if !l.inner.Enabled(ctx, level) {
		return l
	}
	frame := callerFrameForRecord()
	attrs = l.prependModule(frame, attrs)
	attrs = prependSource(frame, attrs)
	r := slog.NewRecord(time.Now(), level, msg, 0)
	r.Add(attrs...)
	_ = l.inner.Handler().Handle(ctx, r)
	return l
}

func (l *Logger) Debug(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelError, msg, args...)
}

func (l *Logger) Debugf(format string, args ...any) *Logger {
	return l.handle(slog.LevelDebug, formatMessage(format, args...), l.attrs)
}

func (l *Logger) Infof(format string, args ...any) *Logger {
	return l.handle(slog.LevelInfo, formatMessage(format, args...), l.attrs)
}

func (l *Logger) Warnf(format string, args ...any) *Logger {
	return l.handle(slog.LevelWarn, formatMessage(format, args...), l.attrs)
}

func (l *Logger) Errorf(format string, args ...any) *Logger {
	return l.handle(slog.LevelError, formatMessage(format, args...), l.attrs)
}

func (l *Logger) emitArgs(level slog.Level, msg string, args ...any) *Logger {
	if len(args) == 0 {
		return l.handle(level, msg, l.attrs)
	}
	if l == nil {
		l = rootLogger()
	}
	attrs := append(append([]any{}, l.attrs...), args...)
	return l.handle(level, msg, attrs)
}

// Package-level helpers (slog-native: msg + optional key-value pairs).

func Debug(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelDebug, msg, args...)
}

func Info(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelError, msg, args...)
}

func Debugf(format string, args ...any) *Logger {
	return rootLogger().Debugf(format, args...)
}

func Infof(format string, args ...any) *Logger {
	return rootLogger().Infof(format, args...)
}

func Warnf(format string, args ...any) *Logger {
	return rootLogger().Warnf(format, args...)
}

func Errorf(format string, args ...any) *Logger {
	return rootLogger().Errorf(format, args...)
}
