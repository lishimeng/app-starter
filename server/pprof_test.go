package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func TestStartPprofSetup(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	pprof.Register(engine.Group(""), DefaultPprofPath)
	engine.GET(DefaultMetricsPath, gin.WrapH(MetricsHandler()))
	called := false
	setup := AdminSetup(func(r Router) {
		called = true
		r.Get("/ping", func(ctx Context) {
			_, _ = ctx.Write([]byte("pong"))
		})
	})
	setup(NewRouter(engine))
	if !called {
		t.Fatal("Setup not called")
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "pong" {
		t.Fatalf("GET /ping = %d %q", w.Code, w.Body.String())
	}
}


func TestAdminRoutes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	pprof.Register(engine.Group(""), DefaultPprofPath)
	engine.GET(DefaultMetricsPath, gin.WrapH(MetricsHandler()))

	req := httptest.NewRequest(http.MethodGet, DefaultPprofPath+"/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("pprof status = %d, want 200", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "go_goroutines") {
		t.Fatalf("metrics body missing go runtime metrics")
	}
}

func TestWatchShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	srv.Start()

	go watchShutdown(ctx, srv.Config)

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	cancel()
	time.Sleep(50 * time.Millisecond)

	_, err = http.Get(srv.URL)
	if err == nil {
		t.Fatal("expected server stopped after ctx cancel")
	}
}

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(MetricsMiddleware())
	engine.GET("/items/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items/42", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	body := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	MetricsHandler().ServeHTTP(body, req)
	if !strings.Contains(body.Body.String(), `route="/items/:id"`) {
		t.Fatalf("metrics missing route label: %s", body.Body.String())
	}
}
