package static_test

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/examples/web-basic/static"
	"github.com/lishimeng/app-starter/server"
)

func TestEmbedStaticFiles(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"index.html", "css/index.css", "js/index.js"} {
		f, err := static.Static.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		_ = f.Close()
	}
}

func TestStaticNoRoute(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/", func(c *gin.Context) { c.String(http.StatusOK, "api") })
	r.NoRoute(server.GinHandler(server.FileServerNoRoute(http.FS(static.Static))))

	for _, path := range []string{"/", "/css/index.css", "/js/index.js"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: status %d", path, rec.Code)
		}
		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatalf("read body %s: %v", path, err)
		}
		if len(body) == 0 {
			t.Fatalf("GET %s: empty body", path)
		}
	}

	index, err := fs.ReadFile(static.Static, "index.html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(index), "web-basic static") {
		t.Fatal("index.html content unexpected")
	}
}
