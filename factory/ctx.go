package factory

import (
	"context"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/cache"
	persistence "github.com/lishimeng/go-orm"
)

var globalContext context.Context
var appCache cache.C

var amqpSession rabbit.Session

func RegisterCtx(ctx context.Context) {
	globalContext = ctx
}

func RegisterCache(c cache.C) {
	appCache = c
}

func RegisterAmqp(session rabbit.Session) {
	amqpSession = session
}

func GetCtx() (ctx context.Context) {
	ctx = globalContext
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
