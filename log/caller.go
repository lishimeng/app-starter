package log

import (
	"path/filepath"
	"runtime"
	"strings"
)

const moduleRoot = "app-starter/"

// moduleFromCaller returns package path relative to app-starter from runtime.Caller(skip).
func moduleFromCaller(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok || file == "" {
		return ""
	}
	file = filepath.ToSlash(file)
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
