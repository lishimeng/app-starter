package persistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	applog "github.com/lishimeng/app-starter/log"
	gormlogger "gorm.io/gorm/logger"
)

// slogGormLogger integrates GORM with app-starter slog.
type slogGormLogger struct {
	logger                    *slog.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func newSlogGormLogger(level gormlogger.LogLevel) gormlogger.Interface {
	return &slogGormLogger{
		logger:                    applog.Slog("gorm"),
		logLevel:                  level,
		slowThreshold:             200 * time.Millisecond,
		ignoreRecordNotFoundError: true,
	}
}

func (l *slogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	next := *l
	next.logLevel = level
	return &next
}

func (l *slogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.log(ctx, slog.LevelInfo, msg, slog.Any("data", data))
	}
}

func (l *slogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.log(ctx, slog.LevelWarn, msg, slog.Any("data", data))
	}
}

func (l *slogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.log(ctx, slog.LevelError, msg, slog.Any("data", data))
	}
}

func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []slog.Attr{
		slog.String("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
		slog.String("sql", sql),
	}
	if rows != -1 {
		fields = append(fields, slog.Int64("rows", rows))
	}

	switch {
	case err != nil && (!l.ignoreRecordNotFoundError || !errors.Is(err, gormlogger.ErrRecordNotFound)):
		fields = append(fields, slog.String("error", err.Error()))
		l.log(ctx, slog.LevelError, "SQL executed", slog.Attr{
			Key: "trace", Value: slog.GroupValue(fields...),
		})
	case l.slowThreshold != 0 && elapsed > l.slowThreshold:
		l.log(ctx, slog.LevelWarn, "SQL executed", slog.Attr{
			Key: "trace", Value: slog.GroupValue(fields...),
		})
	case l.logLevel >= gormlogger.Info:
		l.log(ctx, slog.LevelInfo, "SQL executed", slog.Attr{
			Key: "trace", Value: slog.GroupValue(fields...),
		})
	}
}

func (l *slogGormLogger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.logger.Enabled(ctx, level) {
		return
	}

	frame := gormLogSourceFrame()
	r := slog.NewRecord(time.Now(), level, msg, 0)
	args = applog.PrependSource(frame, args)
	r.Add(args...)
	_ = l.logger.Handler().Handle(ctx, r)
}

// gormLogSourceFrame skips gorm.io internals and this adapter (same idea as gorm.io/gorm/utils.CallerFrame).
func gormLogSourceFrame() runtime.Frame {
	pcs := [32]uintptr{}
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.PC == 0 {
			break
		}
		if gormLogSourceSkip(frame.File) {
			if !more {
				break
			}
			continue
		}
		return frame
	}
	return runtime.Frame{}
}

func gormLogSourceSkip(file string) bool {
	file = filepath.ToSlash(file)
	if strings.Contains(file, "gorm.io/") {
		return true
	}
	if strings.HasSuffix(file, "/persistence/gorm_logger.go") {
		return true
	}
	return false
}

func gormLogLevel(debug bool) gormlogger.LogLevel {
	if debug {
		return gormlogger.Info
	}
	return gormlogger.Silent
}
