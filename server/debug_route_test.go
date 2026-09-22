package server

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/server/pprof"
)

func orgCars(ctx Context) {}

func TestDebugPrintBusinessRouteName(t *testing.T) {
	prevMode := gin.Mode()
	prevWriter := gin.DefaultWriter
	t.Cleanup(func() {
		gin.SetMode(prevMode)
		gin.DefaultWriter = prevWriter
	})
	gin.SetMode(gin.DebugMode)
	var buf bytes.Buffer
	gin.DefaultWriter = &buf

	engine := gin.New()
	engine.Use(gin.Recovery())
	NewRouter(engine).Get("/api/v1/org-cars", orgCars)

	out := buf.String()
	if !strings.Contains(out, "server.orgCars") {
		t.Fatalf("debug log missing real handler name:\n%s", out)
	}
	if strings.Contains(out, "GinHandler.func1") {
		t.Fatalf("debug log still uses adapter name:\n%s", out)
	}
}

func TestDebugPrintPprofNameUnchanged(t *testing.T) {
	prevMode := gin.Mode()
	prevWriter := gin.DefaultWriter
	t.Cleanup(func() {
		gin.SetMode(prevMode)
		gin.DefaultWriter = prevWriter
	})
	gin.SetMode(gin.DebugMode)
	var buf bytes.Buffer
	gin.DefaultWriter = &buf

	engine := gin.New()
	engine.Use(gin.Recovery())
	pprof.Register(engine.Group(""), DefaultPprofPath)

	out := buf.String()
	if !strings.Contains(out, "server/pprof.Index") {
		t.Fatalf("pprof debug log missing named handler:\n%s", out)
	}
	if strings.Contains(out, "GinHandler.func1") {
		t.Fatalf("pprof debug log used adapter name:\n%s", out)
	}
}
