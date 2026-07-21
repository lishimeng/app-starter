package server

import (
	"context"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/log"
)

// DefaultPprofListen is the default admin listen address (pprof + metrics).
const DefaultPprofListen = ":6060"

// DefaultPprofPath is the URL prefix for pprof routes on the admin listener.
const DefaultPprofPath = "/pprof"

// DefaultAdminLogLevelPath is the admin listener path for runtime log level changes.
const DefaultAdminLogLevelPath = "/cl"

// AdminSetup registers custom routes on the admin HTTP listener (:6060).
// Same signature as Component so business code uses server.Router, not gin.
type AdminSetup func(Router)

// PprofConfig configures the admin HTTP listener (pprof, metrics, and admin API).
type PprofConfig struct {
	Listen      string // e.g. DefaultPprofListen; empty disables the listener
	Path        string // pprof prefix, default DefaultPprofPath
	MetricsPath string // default DefaultMetricsPath
	Setup       AdminSetup
}

// StartPprof listens on cfg.Listen and serves pprof and /metrics.
// Returns immediately with nil when Listen is empty.
func StartPprof(ctx context.Context, cfg PprofConfig) error {
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

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	pprof.Register(engine.Group(""), pprofPath)
	engine.GET(metricsPath, gin.WrapH(MetricsHandler()))
	if cfg.Setup != nil {
		cfg.Setup(NewRouter(engine))
	}

	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: engine,
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
