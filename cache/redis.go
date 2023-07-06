package cache

import (
	"context"
	"github.com/go-redis/cache/v9"
	"time"
)

type C interface {
	Get(key string, value interface{}) (err error)
	GetProxy() *cache.Cache
	Exists(key string) bool
	Set(key string, value interface{}) (err error)
	SetTTL(key string, value interface{}, ttl time.Duration) (err error)
	Del(key string) error
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

func (c *redisCache) Exists(key string) bool {
	return c.proxy.Exists(c.ctx, key)
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
	err = c.SetTTL(key, value, 0)
	return
}

func (c *redisCache) SetTTL(key string, value interface{}, ttl time.Duration) (err error) {
	if ttl <= 0 {
		ttl = c.defaultTtl
	}
	item := cache.Item{
		Ctx:   c.ctx,
		Key:   key,
		Value: value,
	}
	if ttl > 0 {
		item.TTL = ttl
	}
	err = c.proxy.Set(&item)
	return
}

// Del 删除, 慎用
func (c *redisCache) Del(key string) (err error) {
	err = c.proxy.Delete(c.ctx, key)
	return
}
