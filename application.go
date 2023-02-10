package app

import (
	"context"
	"fmt"
	"github.com/lishimeng/app-starter/amqp"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/application/repo"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/server"
	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/lishimeng/go-orm"
)

type Application interface {
	Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate func(string)) error
}

type application struct {
	builder *ApplicationBuilder
}

var ctx context.Context
var appCache cache.C

var amqpSession rabbit.Session

func New() (instance Application) {
	ctx = shutdown.Context()
	builder := &ApplicationBuilder{}
	ins := &application{builder: builder}
	instance = ins
	return
}

func GetAmqp() (session rabbit.Session) {
	session = amqpSession
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}

func GetCache() (c cache.C) {
	c = appCache
	return
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
	err = buildHandler(ctx, h.builder)
	if err != nil {
		return
	}

	if h.builder.amqpEnable {
		amqpSession = amqp.New(ctx, h.builder.amqpOptions)
	}

	if h.builder.dbEnable {
		err = repo.Database(h.builder.dbConfig, h.builder.dbModels...)
		if err != nil {
			return
		}
	}

	if h.builder.cacheEnable {
		appCache = cache.New(ctx, h.builder.redisOpts, h.builder.cacheOpts)
	}

	err = h.applyComponents(h.builder.componentsBeforeWebServer)
	if err != nil {
		return err
	}

	if h.builder.webEnable {
		var srv *server.Server
		conf := server.Config{
			Listen: h.builder.webListen,
		}
		if len(h.builder.webLogLevel) > 0 {
			conf.LogLvl = h.builder.webLogLevel
		}
		srv, err = api.Server(conf)
		if h.builder.webStaticEnable {
			err = api.EnableStatic(srv,
				h.builder.assetFile)
			if err != nil {
				return
			}
		}
		err = api.EnableComponents(srv, h.builder.webComponents...)
		if err != nil {
			return
		}
		err = api.EnableMonitors(srv)
		if err != nil {
			return
		}
		err = api.Start(ctx, srv)
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
