package pprof

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterHandlerNames(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	Register(engine, "/pprof")

	wantSuffix := map[string]string{
		"GET /pprof/":             ".Index",
		"GET /pprof/cmdline":      ".Cmdline",
		"GET /pprof/profile":      ".Profile",
		"POST /pprof/symbol":      ".Symbol",
		"GET /pprof/symbol":       ".Symbol",
		"GET /pprof/trace":        ".Trace",
		"GET /pprof/allocs":       ".Allocs",
		"GET /pprof/block":        ".Block",
		"GET /pprof/goroutine":    ".Goroutine",
		"GET /pprof/heap":         ".Heap",
		"GET /pprof/mutex":        ".Mutex",
		"GET /pprof/threadcreate": ".Threadcreate",
	}

	got := map[string]string{}
	for _, ri := range engine.Routes() {
		key := ri.Method + " " + ri.Path
		got[key] = ri.Handler
		if strings.Contains(ri.Handler, "WrapF") || strings.Contains(ri.Handler, "WrapH") {
			t.Errorf("%s handler %q still uses gin.WrapF/WrapH", key, ri.Handler)
		}
	}
	for key, suffix := range wantSuffix {
		name, ok := got[key]
		if !ok {
			t.Errorf("missing route %s", key)
			continue
		}
		if !strings.HasSuffix(name, suffix) {
			t.Errorf("%s handler = %q, want suffix %q", key, name, suffix)
		}
		if !strings.Contains(name, "github.com/lishimeng/app-starter/server/pprof.") {
			t.Errorf("%s handler = %q, want package server/pprof", key, name)
		}
	}
}

func TestIndexOK(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	Register(engine, "/pprof")

	req := httptest.NewRequest(http.MethodGet, "/pprof/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}
