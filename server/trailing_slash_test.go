package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTrimTrailingSlashes(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"/":        "/",
		"":         "",
		"/api":     "/api",
		"/api/":    "/api",
		"/api//":   "/api",
		"/a/b/":    "/a/b",
		"/a/b///":  "/a/b",
	}
	for in, want := range cases {
		if got := trimTrailingSlashes(in); got != want {
			t.Fatalf("trimTrailingSlashes(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestStripTrailingSlashBeforeRouting(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "ok:"+c.Request.URL.Path)
	})

	h := stripTrailingSlash(engine)

	for _, path := range []string{"/api", "/api/", "/api//"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: status %d", path, rec.Code)
		}
		if rec.Body.String() != "ok:/api" {
			t.Fatalf("GET %s: body %q", path, rec.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET / should not be rewritten away, status %d", rec.Code)
	}
}
