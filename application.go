package app

import (
	"context"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/go-libs/etc"
	"github.com/lishimeng/go-libs/persistence"
	server "github.com/lishimeng/go-libs/web"
)

type Application struct {
	_ctx context.Context

	webEnable bool
	webListen string
	webComponents []server.Component

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
		err = api.Server(a._ctx, a.webListen, a.webComponents...)
	}
	if err != nil {
		return
	}

	return
}