package factory

import (
	"context"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/mqtt"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/go-log"
	proxy "github.com/lishimeng/x/container"
)

const (
	amqpKey  = "amqp_session"
	cacheKey = "cache_redis"
	mqttKey  = "mqtt_redis"
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
	proxy.Add(&c, cacheKey)
}

func GetCache() (c cache.C) {
	err := proxy.Get(&c, cacheKey)
	if err != nil {
		log.Debug(err)
		c = nil
	}
	return
}

func RegisterAmqp(session rabbit.Session) {
	proxy.Add(&session, amqpKey)
}

func GetAmqp() (session rabbit.Session) {
	err := proxy.Get(&session, amqpKey)
	if err != nil {
		log.Debug(err)
		session = nil
	}
	return
}

func RegisterMqtt(session mqtt.Session) {
	proxy.Add(&session, mqttKey)
}
func GetMqtt() (session mqtt.Session) {
	err := proxy.Get(&session, mqttKey)
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
