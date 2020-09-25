package app

import (
	"context"
	"github.com/lishimeng/go-etc"
	persistence "github.com/lishimeng/go-orm"
	server "github.com/lishimeng/go-web-server"
)

type ApplicationBuilder struct {

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

func (h *ApplicationBuilder) LoadConfig(config interface{}, name string, path ...string) error {
	_, err := etc.LoadEnvs(name, path, config)
	return err
}

func (h *ApplicationBuilder) EnableWeb(listen string, components ...server.Component) *ApplicationBuilder {
	h.webEnable = true
	h.webListen = listen
	h.webComponents = components
	// TODO check
	return h
}

func (h *ApplicationBuilder) EnableStaticWeb(home string,
	asset func(string) ([]byte, error),
	assetNames func()[]string) *ApplicationBuilder {
	h.webStaticEnable = true
	h.webStaticHome = home
	h.webStaticAsset = asset
	h.webStaticAssetNames = assetNames
	// TODO check
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