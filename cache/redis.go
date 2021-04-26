package cache

import (
	"context"
	"github.com/go-redis/cache/v8"
	"time"
)

type C interface {
	Get(key string, value interface{}) (err error)
	GetProxy() *cache.Cache
	Set(key string, value interface{}) (err error)
}

type redisCache struct {
	ctx        context.Context
	proxy      *cache.Cache
	defaultTtl time.Duration
}

func (c *redisCache) GetProxy() *cache.Cache {
	return c.proxy
}

// Get
/**
err = c.Get(key, &obj)
 */
func (c *redisCache) Get(key string, value interface{}) (err error) {
	err = c.proxy.Get(c.ctx, key, value)
	return
}

// Set
/**
err = c.Set(key, &Demo{
		A: "this is Demo",
		B: 24,
		C: time.Now(),
	})
 */
func (c *redisCache) Set(key string, value interface{}) (err error) {
	err = c.proxy.Set(&cache.Item{
		Ctx:   nil,
		Key:   key,
		Value: value,
		TTL:   c.defaultTtl,
	})
	return
}
