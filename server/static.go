package server

import (
	"io"
	"net/http"
	"path"
	"strings"
)

// FileServerNoRoute serves files from fs for unmatched GET/HEAD requests.
// Use as Gin NoRoute handler so root StaticFS does not conflict with /api, /m, etc.
func FileServerNoRoute(fs http.FileSystem) Handler {
	return func(ctx Context) {
		req := ctx.Request()
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			ctx.Status(http.StatusMethodNotAllowed)
			return
		}
		name := strings.TrimPrefix(req.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}
		name = path.Clean(name)
		if name == "." || strings.HasPrefix(name, "..") || strings.Contains(name, "/..") {
			ctx.Status(http.StatusNotFound)
			return
		}
		serveFileFromFS(ctx, fs, name)
	}
}

func serveFileFromFS(ctx Context, fs http.FileSystem, name string) {
	f, err := fs.Open(name)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		ctx.Status(http.StatusNotFound)
		return
	}

	rs, ok := f.(io.ReadSeeker)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	http.ServeContent(ctx.ResponseWriter(), ctx.Request(), path.Base(name), stat.ModTime(), rs)
}
