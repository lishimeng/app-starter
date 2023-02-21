package factory

import (
	"context"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/go-log"
	persistence "github.com/lishimeng/go-orm"
)

const (
	amqpKey  = "amqp_session"
	cacheKey = "cache_redis"
)

var globalContext context.Context

//var appCache cache.C

//var amqpSession rabbit.Session

func RegisterCtx(ctx context.Context) {
	globalContext = ctx
}

func GetCtx() (ctx context.Context) {
	ctx = globalContext
	return
}

func RegisterCache(c cache.C) {
	Add(&c, cacheKey)
}

func GetCache() (c cache.C) {
	err := Get(&c, cacheKey)
	if err != nil {
		log.Debug(err)
		c = nil
	}
	return
}

func RegisterAmqp(session rabbit.Session) {
	Add(&session, amqpKey)
}

func GetAmqp() (session rabbit.Session) {
	err := Get(&session, amqpKey)
	if err != nil {
		log.Debug(err)
		session = nil
	}
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}
