package app

import (
	"context"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/app-starter/etc"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-orm"
)

type Application interface {
	Etc(config interface{}, callback func(etc.Loader)) error
	EnableWeb(listen string, components ...server.Component) Application
	EnableStaticWeb(home string, asset func(string) ([]byte, error), assetNames func()[]string) Application
	EnableDatabase(config persistence.BaseConfig, models ...interface{}) Application

	ComponentBefore(component func(context.Context)(err error)) Application
	ComponentAfter(component func(context.Context)(err error)) Application

	Start() (err error)
}

type application struct {
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

func New(ctx context.Context) (instance Application) {
	ins := &application{_ctx:ctx}
	instance = ins
	return
}

func GetOrm() *persistence.OrmContext {
	return orm
}

func (h *application) Etc(config interface{}, callback func(etc.Loader)) error {
	var err error
	loader := etc.New()
	if callback != nil {
		callback(loader)
	} else {
		callback = func(ld etc.Loader) {
			ld.SetFileSearcher("config").SetEnvPrefix("").SetEnvSearcher()
		}
	}
	err = loader.Load(config)
	return err
}

func (h *application) EnableWeb(listen string, components ...server.Component) Application {
	h.webEnable = true
	h.webListen = listen
	h.webComponents = components
	// TODO check
	return h
}

func (h *application) EnableStaticWeb(home string, asset func(string) ([]byte, error), assetNames func()[]string) Application {
	h.webStaticEnable = true
	h.webStaticHome = home
	h.webStaticAsset = asset
	h.webStaticAssetNames = assetNames
	// TODO check
	return h
}

func (h *application) EnableDatabase(config persistence.BaseConfig, models ...interface{}) Application {

	h.dbEnable = true
	h.dbConfig = config
	h.dbModels = models
	// TODO check
	return h
}

func (h *application) Start() (err error) {

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