package log

import (
	"context"
	"log/slog"
)

const callerSkip = 5 // package / (*Logger) method -> emitArgs -> resolveModule -> moduleFromCaller

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

func (l *Logger) resolveModule(skip int) string {
	if l != nil && l.module != "" {
		return l.module
	}
	return moduleFromCaller(skip)
}

func (l *Logger) emit(level slog.Level, skip int, msg string) *Logger {
	if l == nil {
		l = rootLogger()
	}
	attrs := l.attrs
	if mod := l.resolveModule(skip); mod != "" {
		attrs = append([]any{slog.String("module", mod)}, attrs...)
	}
	l.inner.Log(context.Background(), level, msg, attrs...)
	return l
}

func (l *Logger) Debug(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelDebug, callerSkip, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelInfo, callerSkip, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelWarn, callerSkip, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) *Logger {
	return l.emitArgs(slog.LevelError, callerSkip, msg, args...)
}

func (l *Logger) Debugf(format string, args ...any) *Logger {
	return l.emit(slog.LevelDebug, callerSkip, formatMessage(format, args...))
}

func (l *Logger) Infof(format string, args ...any) *Logger {
	return l.emit(slog.LevelInfo, callerSkip, formatMessage(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) *Logger {
	return l.emit(slog.LevelWarn, callerSkip, formatMessage(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) *Logger {
	return l.emit(slog.LevelError, callerSkip, formatMessage(format, args...))
}

func (l *Logger) emitArgs(level slog.Level, skip int, msg string, args ...any) *Logger {
	if len(args) == 0 {
		return l.emit(level, skip, msg)
	}
	attrs := append([]any{}, l.attrs...)
	attrs = append(attrs, args...)
	if l == nil {
		l = rootLogger()
	}
	mod := l.resolveModule(skip)
	if mod != "" {
		attrs = append([]any{slog.String("module", mod)}, attrs...)
	}
	l.inner.Log(context.Background(), level, msg, attrs...)
	return l
}

// Package-level helpers (slog-native: msg + optional key-value pairs).

func Debug(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelDebug, callerSkip, msg, args...)
}

func Info(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelInfo, callerSkip, msg, args...)
}

func Warn(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelWarn, callerSkip, msg, args...)
}

func Error(msg string, args ...any) *Logger {
	return rootLogger().emitArgs(slog.LevelError, callerSkip, msg, args...)
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
