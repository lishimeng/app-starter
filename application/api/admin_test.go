package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/server"
)

func TestComposeAdminSetup(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	var order []string
	composeAdminSetup(func(r server.Router) {
		order = append(order, "user")
		r.Get("/custom", func(ctx server.Context) {
			_, _ = ctx.Write([]byte("ok"))
		})
	})(server.NewRouter(engine))

	// Framework /cl must exist even when business Setup is provided.
	req := httptest.NewRequest(http.MethodPost, server.DefaultAdminLogLevelPath, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatal("POST /cl should remain registered (framework API not overwritten)")
	}

	req = httptest.NewRequest(http.MethodGet, "/custom", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "ok" {
		t.Fatalf("GET /custom = %d %q", w.Code, w.Body.String())
	}
	if len(order) != 1 || order[0] != "user" {
		t.Fatalf("user setup order = %v, want [user]", order)
	}
}

func TestComposeAdminSetupOrder(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	var order []string

	// Intercept by wrapping: framework registerAdminRoutes has no hook,
	// so verify /cl is present before user route would conflict if registered first.
	setup := composeAdminSetup(func(r server.Router) {
		order = append(order, "user")
		// If framework /cl were missing, this would be the only POST /cl.
		r.Get("/after-cl", func(ctx server.Context) {
			_, _ = ctx.Write([]byte("after"))
		})
	})
	setup(server.NewRouter(engine))

	if len(order) != 1 || order[0] != "user" {
		t.Fatalf("order = %v", order)
	}

	req := httptest.NewRequest(http.MethodPost, server.DefaultAdminLogLevelPath, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatal("framework POST /cl missing after compose")
	}

	req = httptest.NewRequest(http.MethodGet, "/after-cl", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "after" {
		t.Fatalf("business route missing: %d %q", w.Code, w.Body.String())
	}
}

func TestComposeAdminSetupNilUser(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	composeAdminSetup(nil)(server.NewRouter(engine))

	req := httptest.NewRequest(http.MethodPost, server.DefaultAdminLogLevelPath, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatal("POST /cl should be registered without user setup")
	}
}
