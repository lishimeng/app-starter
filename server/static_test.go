package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestFileServerNoRoute(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/api/", func(c *gin.Context) { c.String(http.StatusOK, "api") })
	r.GET("/m/healthy", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.NoRoute(wrapGinHandler(FileServerNoRoute(http.Dir("testdata/static"))))

	for _, tc := range []struct {
		path   string
		status int
		body   string
	}{
		{"/api/", http.StatusOK, "api"},
		{"/m/healthy", http.StatusOK, "ok"},
		{"/", http.StatusOK, "web-basic static"},
		{"/css/index.css", http.StatusOK, "aquamarine"},
		{"/js/index.js", http.StatusOK, "static assets loaded"},
		{"/missing.txt", http.StatusNotFound, ""},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Fatalf("GET %s: status %d, want %d", tc.path, rec.Code, tc.status)
		}
		if tc.body != "" && !strings.Contains(rec.Body.String(), tc.body) {
			t.Fatalf("GET %s: body %q, want substring %q", tc.path, rec.Body.String(), tc.body)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /: status %d, want 405", rec.Code)
	}
}

func TestFileServerNoRouteWithAPIRegisteredAfter(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.NoRoute(wrapGinHandler(FileServerNoRoute(http.Dir("testdata/static"))))
	r.GET("/api/", func(c *gin.Context) { c.String(http.StatusOK, "api") })

	req := httptest.NewRequest(http.MethodGet, "/api/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "api") {
		t.Fatalf("GET /api/: %d %s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	body, _ := io.ReadAll(rec.Body)
	if rec.Code != http.StatusOK || !strings.Contains(string(body), "web-basic static") {
		t.Fatalf("GET /: %d %s", rec.Code, body)
	}
}
