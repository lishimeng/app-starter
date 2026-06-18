package log

import (
	"fmt"
	"log/slog"
)

// LevelFromString maps go-log style level names to slog.Level.
func LevelFromString(s string) (slog.Level, error) {
	switch stringsUpperPrefix(s) {
	case "FINEST", "FINE", "DEBUG":
		return slog.LevelDebug, nil
	case "TRACE", "INFO":
		return slog.LevelInfo, nil
	case "WARNING", "WARN":
		return slog.LevelWarn, nil
	case "ERROR", "CRITICAL":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("log: unknown level %q", s)
	}
}

func stringsUpperPrefix(s string) string {
	if s == "" {
		return ""
	}
	end := 0
	for end < len(s) && s[end] != ' ' && s[end] != '\t' {
		end++
	}
	return stringsToUpper(s[:end])
}

func stringsToUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// formatMessage is used by *f helpers (printf-style only).
func formatMessage(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
