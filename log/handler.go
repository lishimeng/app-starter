package log

import (
	"io"
	"log/slog"
)

func newHandler(w io.Writer, level *slog.LevelVar, json bool) slog.Handler {
	opts := &slog.HandlerOptions{Level: level}
	if json {
		return slog.NewJSONHandler(w, opts)
	}
	return slog.NewTextHandler(w, opts)
}
