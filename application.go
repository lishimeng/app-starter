package app

import (
	"context"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/mqtt"
	"github.com/lishimeng/app-starter/persistence"
)

type Application interface {
	Start(buildHandler func(ctx context.Context, builder *ApplicationBuilder) error, onTerminate ...func(string)) error
}

func GetWebServer() {

}

func GetAmqp() (session rabbit.Session) {
	session = factory.GetAmqp()
	return
}

func GetMqtt() (session mqtt.Session) {
	session = factory.GetMqtt()
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}

func GetCache() (c cache.C) {
	c = factory.GetCache()
	return
}
