package app

import (
	"context"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/etc"
	"github.com/lishimeng/app-starter/server"
	persistence "github.com/lishimeng/go-orm"
	"os"
)

const (
	defaultWebLogLevel = "error"
)

type ApplicationBuilder struct {
	webEnable     bool
	webListen     string
	webComponents []server.Component

	webStaticEnable bool
	vdir            string
	assetInfo       func(name string) (os.FileInfo, error)
	asset           func(string) ([]byte, error)
	assetNames      func() []string
	webStaticHome   string

	webLogLevel string

	dbEnable bool
	dbConfig persistence.BaseConfig
	dbModels []interface{}

	cacheEnable bool
	redisOpts   cache.RedisOptions
	cacheOpts   cache.Options

	// other components
	componentsBeforeWebServer []func(ctx context.Context) (err error)
	componentsAfterWebServer  []func(ctx context.Context) (err error)
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
	h.webComponents = append(h.webComponents, api.Router)
	if len(components) > 0 {
		h.webComponents = append(h.webComponents, components...)
	}
	// TODO check
	return h
}

func (h *ApplicationBuilder) SetMonitorPrefix(prefix string) *ApplicationBuilder {
	api.MonitorPrefix = prefix
	return h
}

func (h *ApplicationBuilder) HealthyHandler(handler func() int) *ApplicationBuilder {
	if handler != nil {
		api.LivenessHandler = handler
	}
	return h
}

func (h *ApplicationBuilder) ReadyHandler(handler func() int) *ApplicationBuilder {
	if handler != nil {
		api.ReadinessHandler = handler
	}
	return h
}

func (h *ApplicationBuilder) SetWebLogLevel(lvl string) *ApplicationBuilder {
	h.webLogLevel = lvl
	return h
}

func (h *ApplicationBuilder) EnableStaticWeb(vdir, home string,
	assetInfo func(name string) (os.FileInfo, error),
	asset func(string) ([]byte, error),
	assetNames func() []string) *ApplicationBuilder {
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

func (h *ApplicationBuilder) EnableCache(redisOpts cache.RedisOptions, cacheOpts cache.Options) *ApplicationBuilder {

	h.cacheEnable = true
	h.redisOpts = redisOpts
	h.cacheOpts = cacheOpts
	return h
}