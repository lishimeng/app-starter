package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/log"
)

// DefaultAdminListen is the default admin listen address (pprof + metrics + admin API).
const DefaultAdminListen = ":6060"

// DefaultPprofPath is the URL prefix for pprof routes on the admin listener.
const DefaultPprofPath = "/pprof"

// DefaultAdminLogLevelPath is the admin listener path for runtime log level changes.
const DefaultAdminLogLevelPath = "/cl"

// AdminSetup registers custom routes on the admin HTTP listener (:6060).
// Same signature as Component so business code uses server.Router, not gin.
type AdminSetup func(Router)

// AdminConfig configures the admin HTTP listener (pprof, metrics, and admin API).
type AdminConfig struct {
	Listen             string // e.g. DefaultAdminListen; empty disables the listener
	Path               string // pprof prefix, default DefaultPprofPath
	MetricsPath        string // default DefaultMetricsPath
	LogLvl             string // "debug" enables Gin route listing (same as web LogLvl)
	StripTrailingSlash bool   // rewrite /demo/ping/ → /demo/ping before routing (no 301)
	Setup              AdminSetup
}

// StartAdmin listens on cfg.Listen and serves pprof, /metrics, and admin routes.
// Returns immediately with nil when Listen is empty.
func StartAdmin(ctx context.Context, cfg AdminConfig) error {
	if cfg.Listen == "" {
		return nil
	}
	pprofPath := cfg.Path
	if pprofPath == "" {
		pprofPath = DefaultPprofPath
	}
	metricsPath := cfg.MetricsPath
	if metricsPath == "" {
		metricsPath = DefaultMetricsPath
	}

	if strings.EqualFold(cfg.LogLvl, "debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	if cfg.StripTrailingSlash {
		engine.RedirectTrailingSlash = false
	}
	engine.Use(gin.Recovery())
	pprof.Register(engine.Group(""), pprofPath)
	engine.GET(metricsPath, gin.WrapH(MetricsHandler()))
	if cfg.Setup != nil {
		cfg.Setup(NewRouter(engine))
	}

	handler := http.Handler(engine)
	if cfg.StripTrailingSlash {
		handler = stripTrailingSlash(handler)
	}
	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: handler,
	}
	return serveHTTP(ctx, srv, "admin server")
}

// serveHTTP starts srv and shuts it down when ctx is cancelled.
func serveHTTP(ctx context.Context, srv *http.Server, label string) error {
	go watchShutdown(ctx, srv)
	log.Infof("%s listen %s", label, srv.Addr)
	return srv.ListenAndServe()
}

func watchShutdown(ctx context.Context, srv *http.Server) {
	if ctx == nil {
		return
	}
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}
