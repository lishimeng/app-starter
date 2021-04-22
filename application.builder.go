package app

import (
	"context"
	"github.com/lishimeng/app-starter/etc"
	"github.com/lishimeng/app-starter/server"
	persistence "github.com/lishimeng/go-orm"
	"os"
)

type ApplicationBuilder struct {

	webEnable bool
	webListen string
	webComponents []server.Component

	webStaticEnable bool
	vdir string
	assetInfo  func(name string) (os.FileInfo, error)
	asset func(string) ([]byte, error)
	assetNames func()[]string
	webStaticHome string

	dbEnable bool
	dbConfig persistence.BaseConfig
	dbModels []interface{}

	// other components
	componentsBeforeWebServer []func(ctx context.Context) (err error)
	componentsAfterWebServer []func(ctx context.Context) (err error)
}

func (h *ApplicationBuilder) LoadConfig(config interface{}, callback func(etc.Loader)) (err error) {
	var cb = callback
	loader := etc.New()
	if cb == nil {
		cb = func(ld etc.Loader) {
			ld.SetFileSearcher("config")
		}
	}
	cb(loader)
	err = loader.Load(config)
	return
}

func (h *ApplicationBuilder) EnableWeb(listen string, components ...server.Component) *ApplicationBuilder {
	h.webEnable = true
	h.webListen = listen
	h.webComponents = components
	// TODO check
	return h
}

func (h *ApplicationBuilder) EnableStaticWeb(vdir, home string,
	assetInfo  func(name string) (os.FileInfo, error),
	asset func(string) ([]byte, error),
	assetNames func()[]string) *ApplicationBuilder {
	h.webStaticEnable = true
	h.vdir = vdir
	h.webStaticHome = home
	h.assetInfo = assetInfo
	h.asset = asset
	h.assetNames = assetNames
	return h
}

func (h *ApplicationBuilder) EnableDatabase(config persistence.BaseConfig,
	models ...interface{}) *ApplicationBuilder {

	h.dbEnable = true
	h.dbConfig = config
	h.dbModels = models
	// TODO check
	return h
}