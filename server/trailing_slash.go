package server

import (
	"net/http"
	"strings"
)

// stripTrailingSlash wraps next so paths like /api/ are rewritten to /api
// before routing (nginx-style, no redirect). Root "/" is unchanged.
func stripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if cleaned := trimTrailingSlashes(path); cleaned != path {
			r2 := r.Clone(r.Context())
			r2.URL.Path = cleaned
			if r.URL.RawPath != "" {
				r2.URL.RawPath = trimTrailingSlashes(r.URL.RawPath)
			}
			next.ServeHTTP(w, r2)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func trimTrailingSlashes(path string) string {
	if path == "" || path == "/" {
		return path
	}
	return strings.TrimRight(path, "/")
}
