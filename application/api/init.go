package api

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
)

func Server(listen string) (srv *server.Server, err error) {
	if len(listen) == 0 {
		return
	}

	srv = server.New(server.Config{
		Listen: listen,
	})
	return
}

func EnableComponents(srv *server.Server, components ...server.Component) (err error) {

	if len(components) == 0 {
		return
	}
	srv.RegisterComponents(components...)
	return
}

func EnableStatic(srv *server.Server, home string, asset func(string) ([]byte, error), assetNames func()[]string) (err error) {

	bs, err := asset(home)
	indexHtml := ""
	if err != nil {
		return
	}
	indexHtml = string(bs)
	srv.SetHomePage(indexHtml)
	srv.AdvancedConfig(func(app *iris.Application) {
		app.HandleDir("/", "", iris.DirOptions{
			Asset: asset,
			AssetNames: assetNames,
		})
	})
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