package factory

import (
	"context"

	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/mqtt"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/redis"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/log"
	proxy "github.com/lishimeng/x/container"
)

const (
	cacheKey     = "cache_redis"
	redisKey     = "redis_client"
	mqttKey      = "mqtt_key"
	webServerKey = "webserver_key"
)

var globalContext context.Context

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
		log.Debugf("%v", err)
		c = nil
	}
	return
}

// RegisterRedis creates and registers a shared redis.Client; closes on app ctx done.
func RegisterRedis(opts redis.Options) {
	client := redis.NewClient(opts)
	if client == nil {
		return
	}
	proxy.Add(&client, redisKey)
	go func() {
		ctx := GetCtx()
		select {
		case <-ctx.Done():
			_ = client.Close()
		}
	}()
}

func GetRedisClient() (c *redis.Client) {
	err := proxy.Get(&c, redisKey)
	if err != nil {
		log.Debugf("%v", err)
		c = nil
	}
	return
}

func RegisterMqtt(session mqtt.Session) {
	proxy.Add(&session, mqttKey)
}
func GetMqtt() (session mqtt.Session) {
	err := proxy.Get(&session, mqttKey)
	if err != nil {
		log.Debugf("%v", err)
		session = nil
	}
	return
}

func RegisterWebServer(s server.Server) {
	proxy.Add(&s, webServerKey)
}
func GetWebServer() (s *server.Server) {
	var svr server.Server
	err := proxy.Get(&svr, webServerKey)
	if err != nil {
		log.Debugf("%v", err)
		s = nil
	} else {
		s = &svr
	}
	return
}

func GetOrm() *persistence.OrmContext {
	return persistence.New()
}
func GetNamedOrm(aliaName string) *persistence.OrmContext {
	return persistence.NewOrm(aliaName)
}
