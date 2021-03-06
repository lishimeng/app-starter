package api

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
	"os"
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

func EnableStatic(srv *server.Server, vdir string, home string,
	assetInfo func(name string) (os.FileInfo, error),
	asset func(string) ([]byte, error),
	assetNames func() []string) (err error) {

	bs, err := asset(vdir + "/" + home)
	indexHtml := ""
	if err != nil {
		return
	}
	indexHtml = string(bs)
	srv.SetHomePage(indexHtml)
	srv.AdvancedConfig(func(app *iris.Application) {
		app.HandleDir("/", vdir, iris.DirOptions{
			AssetInfo:  assetInfo,
			Asset:      asset,
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
