package api

import (
	"context"
	"github.com/kataras/iris/v12"

	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
	"net/http"
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

func EnableMonitors(srv *server.Server) (err error) {
	var r = server.NewRouter(srv.GetMonitor())
	Router(r)
	return
}

func EnableStatic(srv *server.Server, assetFile func() http.FileSystem) (err error) {

	srv.AdvancedConfig(func(app *iris.Application) {
		app.HandleDir("/", assetFile())
	})
	return
}

func EnableVue3Plugin(srv *server.Server, handler func(app *iris.Application)) (err error) {
	srv.AdvancedConfig(handler)
	return
}

func Start(ctx context.Context, srv *server.Server) (err error) {
	go func() {
		log.Info("start web server")

		e := srv.Start(ctx)
		if e != nil {
			log.Info(e)
		}
		log.Info("stop web server")
	}()
	return nil
}
