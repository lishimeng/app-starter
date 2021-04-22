package app

import (
	"context"
	"fmt"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/app-starter/server"
	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/lishimeng/go-orm"
)

type Application interface {
	Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate func(string)) error
}

type application struct {
	_ctx context.Context
	builder *ApplicationBuilder
}

var orm *persistence.OrmContext

func New() (instance Application) {
	ctx := shutdown.Context()
	builder := &ApplicationBuilder{}
	ins := &application{_ctx:ctx, builder: builder}
	instance = ins
	return
}

func GetOrm() *persistence.OrmContext {
	return orm
}

func (h *application) Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate func(string)) (err error) {

	err = h._start(buildHandler)

	if err == nil {
		shutdown.WaitExit(&shutdown.Configuration{
			BeforeExit: func(s string) {
				if onTerminate != nil {
					onTerminate(s)
				}
			},
		})
	}
	return
}

func (h *application) _start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error) (err error) {

	if buildHandler == nil {
		err = fmt.Errorf("application builder function nil")
		return
	}
	err = buildHandler(h._ctx, h.builder)
	if err != nil {
		return
	}
	if h.builder.dbEnable {
		orm, err =repo.Database(h.builder.dbConfig, h.builder.dbModels...)
		if err != nil {
			return
		}
	}
	err = h.applyComponents(h.builder.componentsBeforeWebServer)
	if err != nil {
		return err
	}

	if h.builder.webEnable {
		var srv *server.Server
		srv, err = api.Server(h.builder.webListen)
		if h.builder.webStaticEnable {
			err = api.EnableStatic(srv,
				h.builder.webStaticHome,
				h.builder.webStaticAsset,
				h.builder.webStaticAssetNames)
			if err != nil {
				return
			}
		}
		err = api.EnableComponents(srv, h.builder.webComponents...)
		if err != nil {
			return
		}
		err = api.Start(h._ctx, srv)
		if err != nil {
			return
		}
	}

	err = h.applyComponents(h.builder.componentsAfterWebServer)
	if err != nil {
		return err
	}

	return
}