package log

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

const moduleRoot = "app-starter/"

const logFuncPrefix = "github.com/lishimeng/app-starter/log."

// CallerFrame returns the first stack frame outside log package wrappers (business log only).
func CallerFrame() runtime.Frame {
	return walkCallerFrames(1)
}

func callerFrameForRecord() runtime.Frame {
	return walkCallerFrames(2)
}

func walkCallerFrames(skip int) runtime.Frame {
	pcs := [32]uintptr{}
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.PC == 0 {
			break
		}
		if callerFrameSkip(frame) {
			if !more {
				break
			}
			continue
		}
		return frame
	}
	return runtime.Frame{}
}

func callerFrameSkip(frame runtime.Frame) bool {
	if strings.HasPrefix(frame.Function, logFuncPrefix) {
		rest := strings.TrimPrefix(frame.Function, logFuncPrefix)
		if strings.HasPrefix(rest, "Test") {
			return false
		}
		return true
	}
	file := filepath.ToSlash(frame.File)
	if strings.Contains(file, moduleRoot+"log/") {
		switch filepath.Base(file) {
		case "logger.go", "handler.go", "caller.go", "config.go", "format.go":
			return true
		}
	}
	return false
}

func moduleFromFrame(frame runtime.Frame) string {
	if frame.File == "" {
		return ""
	}
	file := filepath.ToSlash(frame.File)
	idx := strings.Index(file, moduleRoot)
	if idx < 0 {
		return ""
	}
	rel := file[idx+len(moduleRoot):]
	dir := filepath.Dir(rel)
	if dir == "." {
		return ""
	}
	return filepath.ToSlash(dir)
}

func relCallerFromFile(file string, line int) string {
	file = filepath.ToSlash(file)
	idx := strings.Index(file, moduleRoot)
	if idx < 0 {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	rel := file[idx+len(moduleRoot):]
	return fmt.Sprintf("%s:%d", rel, line)
}

func relSourceFromFrame(frame runtime.Frame) string {
	if frame.File == "" {
		return ""
	}
	return relCallerFromFile(frame.File, frame.Line)
}

func prependSource(frame runtime.Frame, attrs []any) []any {
	if !sourceEnabled() {
		return attrs
	}
	if src := relSourceFromFrame(frame); src != "" {
		return append([]any{slog.String("source", src)}, attrs...)
	}
	return attrs
}

// PrependSource adds source=file:line when enabled (for slog integrations such as GORM).
func PrependSource(frame runtime.Frame, attrs []any) []any {
	return prependSource(frame, attrs)
}
