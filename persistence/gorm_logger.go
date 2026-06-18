package persistence

import (
	"fmt"
	"time"

	applog "github.com/lishimeng/app-starter/log"
	gormlogger "gorm.io/gorm/logger"
)

// gormSlogWriter forwards GORM's ANSI-colored SQL lines via log.WriteRaw (bypasses slog escaping).
type gormSlogWriter struct{}

func (gormSlogWriter) Printf(format string, args ...interface{}) {
	applog.WriteRaw(fmt.Sprintf(format, args...))
}

func newSlogGormLogger(level gormlogger.LogLevel) gormlogger.Interface {
	return gormlogger.New(gormSlogWriter{}, gormlogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	}).LogMode(level)
}

func gormLogLevel(debug bool) gormlogger.LogLevel {
	if debug {
		return gormlogger.Info
	}
	return gormlogger.Silent
}
