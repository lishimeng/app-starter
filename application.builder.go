package app

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/amqp"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/etc"
	"github.com/lishimeng/app-starter/mqtt"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/app-starter/version"
	"github.com/lishimeng/go-log"
	persistence "github.com/lishimeng/go-orm"
	"net/http"
	"os"
)

type TokenValidatorInjectFunc func(storage token.Storage)
type TokenValidatorBuilder func(injectFunc TokenValidatorInjectFunc)

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
	assetFile       func() http.FileSystem

	webLogLevel string

	dbEnable bool
	dbConfig persistence.BaseConfig
	dbModels []interface{}

	cacheEnable bool
	redisOpts   cache.RedisOptions
	cacheOpts   cache.Options

	tokenValidatorEnable  bool
	tokenValidatorBuilder TokenValidatorBuilder

	amqpEnable     bool
	amqpOptions    amqp.Connector
	amqpHandler    []amqp.Handler
	sessionOptions []rabbit.SessionOption

	mqttEnable  bool
	mqttOptions []mqtt.ClientOption

	// other components
	componentsBeforeWebServer []func(ctx context.Context) (err error)
	componentsAfterWebServer  []func(ctx context.Context) (err error)
}

func (h *ApplicationBuilder) LoadConfig(config interface{}, callback func(etc.Loader)) (err error) {
	var cb = callback
	loader := etc.New()
	if cb == nil {
		cb = func(ld etc.Loader) {
			ld.SetFileSearcher("config", ".").SetEnvPrefix("").SetEnvSearcher()
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

func (h *ApplicationBuilder) EnableStaticWeb(assetFile func() http.FileSystem) *ApplicationBuilder {
	h.webStaticEnable = true
	h.assetFile = assetFile
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

func (h *ApplicationBuilder) EnableOrmLog() *ApplicationBuilder {
	orm.Debug = true
	return h
}

func (h *ApplicationBuilder) EnableAmqp(c amqp.Connector, options ...rabbit.SessionOption) *ApplicationBuilder {
	h.amqpEnable = true
	h.amqpOptions = c
	h.sessionOptions = append(h.sessionOptions, options...)
	return h
}

// RegisterAmqpHandlers 注册amqp handler
//
// 业务类任务使用延时执行策略，在连接型任务之后执行
func (h *ApplicationBuilder) RegisterAmqpHandlers(handlers ...amqp.Handler) *ApplicationBuilder {
	h.amqpHandler = append(h.amqpHandler, handlers...)
	return h
}

func (h *ApplicationBuilder) EnableMqtt(options ...mqtt.ClientOption) *ApplicationBuilder {
	h.mqttEnable = true
	if len(options) > 0 {
		h.mqttOptions = append(h.mqttOptions, options...)
	}
	log.Debug("enable mqtt module")
	return h
}

// EnableTokenValidator 验证Token，使用RedisTokenValidator前需要enableCache
func (h *ApplicationBuilder) EnableTokenValidator(builder TokenValidatorBuilder) *ApplicationBuilder {
	h.tokenValidatorEnable = true
	h.tokenValidatorBuilder = builder
	return h
}

func (h *ApplicationBuilder) PrintVersion() *ApplicationBuilder {
	version.Print()
	return h
}
