package api

import (
	"context"
	"github.com/kataras/iris"
	"github.com/lishimeng/go-libs/log"
	server "github.com/lishimeng/go-libs/web"
)

func Server(ctx context.Context, listen string, components ...server.Component) (err error) {
	if len(listen) == 0 {
		return
	}
	if len(components) == 0 {
		return
	}

	go func() {
		log.Info("start web server")
		s := server.New(server.ServerConfig{
			Listen: listen,
		})
		s.OnErrorCode(404, func(ctx iris.Context) {

			_, _ = ctx.Text("not found")
		})
		s.RegisterComponents(components...)
		err = s.Start(ctx)
		if err != nil {
			log.Info(err)
		}
		log.Info("stop web server")
	}()
	return
}