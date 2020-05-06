package app

import (
	"context"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/go-etc"
	"github.com/lishimeng/go-orm"
	server "github.com/lishimeng/go-web-server"
)

type Application struct {
	_ctx context.Context

	webEnable bool
	webListen string
	webComponents []server.Component

	webStaticEnable bool
	webStaticAsset func(string) ([]byte, error)
	webStaticAssetNames func()[]string
	webStaticHome string

	dbEnable bool
	dbConfig persistence.BaseConfig
	dbModels []interface{}

	// other components
	componentsBeforeWebServer []func(ctx context.Context) (err error)
	componentsAfterWebServer []func(ctx context.Context) (err error)
}

var orm *persistence.OrmContext

func New(ctx context.Context) (instance *Application) {
	instance = &Application{_ctx:ctx}
	return
}

func GetOrm() *persistence.OrmContext {
	return orm
}

func GetServer() {

}

func (h *Application) LoadConfig(config interface{}, name string, path ...string) error {
	_, err := etc.LoadEnvs(name, path, config)
	return err
}

func (h *Application) EnableWeb(listen string, components ...server.Component) *Application {
	h.webEnable = true
	h.webListen = listen
	h.webComponents = components
	// TODO check
	return h
}

func (h *Application) EnableStaticWeb(home string, asset func(string) ([]byte, error), assetNames func()[]string) *Application {
	h.webStaticEnable = true
	h.webStaticHome = home
	h.webStaticAsset = asset
	h.webStaticAssetNames = assetNames
	// TODO check
	return h
}

func (h *Application) EnableDatabase(config persistence.BaseConfig, models ...interface{}) *Application {

	h.dbEnable = true
	h.dbConfig = config
	h.dbModels = models
	// TODO check
	return h
}

func (h *Application) Start() (err error) {

	if h.dbEnable {
		orm, err =repo.Database(h.dbConfig, h.dbModels...)
		if err != nil {
			return err
		}
	}
	err = h.applyComponents(h.componentsBeforeWebServer)
	if err != nil {
		return err
	}

	if h.webEnable {
		var srv *server.Server
		srv, err = api.Server(h.webListen)
		if h.webStaticEnable {
			err = api.EnableStatic(srv, h.webStaticHome, h.webStaticAsset, h.webStaticAssetNames)
			if err != nil {
				return
			}
		}
		err = api.EnableComponents(srv, h.webComponents...)
		if err != nil {
			return
		}
		err = api.Start(h._ctx, srv)
		if err != nil {
			return
		}
	}

	err = h.applyComponents(h.componentsAfterWebServer)
	if err != nil {
		return err
	}

	return
}