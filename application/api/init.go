package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

func Server(conf server.Config) (srv *server.Server, err error) {
	srv = server.New(conf)
	return
}

func EnableComponents(srv *server.Server, components ...server.Component) (err error) {
	if len(components) == 0 {
		return
	}
	srv.RegisterComponents(components...)
	return
}

func EnableStatic(srv *server.Server, assetFile func() http.FileSystem) (err error) {
	// Gin cannot register StaticFS("/") when /api, /m, etc. exist on the same engine.
	srv.SetNoRoute(server.FileServerNoRoute(assetFile()))
	return
}

func Start(ctx context.Context, srv *server.Server) (err error) {
	go func() {
		log.Info("start web server")
		e := srv.Start(ctx)
		if e != nil {
			log.Infof("%v", e)
		}
		log.Info("stop web server")
	}()
	return nil
}

func StartPprof(ctx context.Context, cfg server.PprofConfig) (err error) {
	cfg.Setup = registerAdminRoutes
	go func() {
		log.Info("start admin server")
		e := server.StartPprof(ctx, cfg)
		if e != nil && e != http.ErrServerClosed {
			log.Infof("%v", e)
		}
		log.Info("stop admin server")
	}()
	return nil
}

func registerAdminRoutes(e *gin.Engine) {
	e.POST(server.DefaultAdminLogLevelPath, server.GinHandler(changeLogLevel))
}
