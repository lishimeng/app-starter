package app

import (
	"context"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/go-libs/etc"
	"github.com/lishimeng/go-libs/persistence"
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
}

var orm *persistence.OrmContext

func New(ctx context.Context) (instance *Application) {
	instance = &Application{_ctx:ctx}
	return
}

func GetOrm() *persistence.OrmContext {
	return orm
}

func (a *Application) LoadConfig(config interface{}, name string, path ...string) error {
	_, err := etc.LoadEnvs(name, path, config)
	return err
}

func (a *Application) EnableWeb(listen string, components ...server.Component) *Application {
	a.webEnable = true
	a.webListen = listen
	a.webComponents = components
	// TODO check
	return a
}

func (a *Application) EnableStaticWeb(home string, asset func(string) ([]byte, error), assetNames func()[]string) *Application {
	a.webStaticEnable = true
	a.webStaticHome = home
	a.webStaticAsset = asset
	a.webStaticAssetNames = assetNames
	// TODO check
	return a
}

func (a *Application) EnableDatabase(config persistence.BaseConfig, models ...interface{}) *Application {

	a.dbEnable = true
	a.dbConfig = config
	a.dbModels = models
	// TODO check
	return a
}

func (a *Application) Start() (err error) {

	if a.dbEnable {
		orm, err =repo.Database(a.dbConfig, a.dbModels...)
	}

	if a.webEnable {
		var srv *server.Server
		srv, err = api.Server(a.webListen)
		if a.webStaticEnable {
			err = api.EnableStatic(srv, a.webStaticHome, a.webStaticAsset, a.webStaticAssetNames)
			if err != nil {
				return
			}
		}
		err = api.EnableComponents(srv, a.webComponents...)
		if err != nil {
			return
		}
		err = api.Start(a._ctx, srv)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	return
}