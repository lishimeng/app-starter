package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	app "github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/web-basic/admin"
	"github.com/lishimeng/app-starter/server"
	shutdown "github.com/lishimeng/go-app-shutdown"
)

func TestRegisterOnEngine(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	admin.Register(server.NewRouter(engine))

	req := httptest.NewRequest(http.MethodGet, admin.DemoPath, nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: status %d", admin.DemoPath, rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != true {
		t.Fatalf("body = %#v", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/demo/health", nil)
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.String() != "admin-ok" {
		t.Fatalf("GET /demo/health: %d %q", rec.Code, rec.Body.String())
	}
}

func TestAdminRoutesWithApp(t *testing.T) {
	adminAddr := freeListenAddr(t)
	webAddr := freeListenAddr(t)

	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(300 * time.Millisecond)

		base := "http://" + adminAddr
		resp, err := http.Get(base + admin.DemoPath)
		if err != nil {
			t.Errorf("GET demo: %v", err)
			shutdown.Exit("test fail")
			return
		}
		raw, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s: status %d body %s", admin.DemoPath, resp.StatusCode, raw)
		}

		resp, err = http.Get(base + "/demo/health")
		if err != nil {
			t.Errorf("GET health: %v", err)
			shutdown.Exit("test fail")
			return
		}
		raw, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK || string(raw) != "admin-ok" {
			t.Errorf("GET /demo/health: %d %q", resp.StatusCode, raw)
		}

		resp, err = http.Get(base + server.DefaultMetricsPath)
		if err != nil {
			t.Errorf("GET metrics: %v", err)
			shutdown.Exit("test fail")
			return
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /metrics: status %d", resp.StatusCode)
		}

		shutdown.Exit("test done")
	}()

	err := app.New().Start(func(ctx context.Context, builder *app.ApplicationBuilder) error {
		builder.
			SetAdminListen(adminAddr).
			EnableAdminRoutes(admin.Register).
			EnableWeb(webAddr)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	<-done
}

func freeListenAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}
