package app

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/server"
)

//go:embed testdata/embed/*
var sample embed.FS

func TestWithEmbedRoot(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	handler := WithEmbedRoot(sample, "testdata/embed")
	r := gin.New()
	r.NoRoute(server.GinHandler(server.FileServerNoRoute(handler())))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /: status %d", rec.Code)
	}
	if !contains(rec.Body.String(), "embed-root") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestWithEmbed(t *testing.T) {
	t.Parallel()
	handler := WithEmbed(sample)
	if handler() == nil {
		t.Fatal("nil filesystem")
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && stringIndex(s, sub) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
