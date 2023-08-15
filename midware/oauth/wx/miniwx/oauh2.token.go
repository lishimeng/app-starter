package miniwx

import (
	"github.com/lishimeng/app-starter/cache"
	"sync"
)

type Oauth2Handle interface {
	// Credentials 直接换取accessToken(服务器对服务器应用)
	Credentials() (accessToken string, err error)
	// AuthorizeCode 微信登录
	AuthorizeCode(code string) (accessToken WxMiniLoginResp, err error)
}

type BaseHandler struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.C
	cacheEnable     bool
	accessTokenLock *sync.Mutex
}

func newBaseHandler() BaseHandler {
	return BaseHandler{
		accessTokenLock: new(sync.Mutex),
	}
}

type Option func(handler *BaseHandler)

func WithAuth(appID string, appSecret string) Option {
	return func(handler *BaseHandler) {
		handler.appID = appID
		handler.appSecret = appSecret
	}
}

func WithCache(c cache.C) Option {
	return func(handler *BaseHandler) {
		handler.cache = c
		handler.cacheEnable = true
	}
}
