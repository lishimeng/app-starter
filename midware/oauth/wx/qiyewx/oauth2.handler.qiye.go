package qiyewx

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/x/rest"
	"sync"
	"time"
)

type AccessTokenHandle interface {
	// GetAccessToken 直接换取accessToken(服务器对服务器应用)
	GetAccessToken() (accessToken string, err error)
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
	}
}

type BaseHandler struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.C
	accessTokenLock *sync.Mutex
}

func newBaseHandler() BaseHandler {
	return BaseHandler{
		accessTokenLock: new(sync.Mutex),
	}
}

// WorkAccessToken 企业微信AccessToken 获取
type WorkAccessToken struct {
	handler BaseHandler
}

// NewWorkAccessToken new WorkAccessToken
func NewWorkAccessToken(opts ...Option) AccessTokenHandle {

	handler := newBaseHandler()
	for _, opt := range opts {
		opt(&handler)
	}
	if handler.cache == nil {
		panic("cache the not exist")
	}
	return &WorkAccessToken{
		handler: handler,
	}
}

// GetAccessToken 企业微信获取access_token,先从cache中获取，没有则从服务端获取
func (ak *WorkAccessToken) GetAccessToken() (accessToken string, err error) {
	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.handler.accessTokenLock.Lock()
	defer ak.handler.accessTokenLock.Unlock()
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.handler.cacheKeyPrefix, ak.handler.appID)
	err = ak.handler.cache.Get(accessTokenCacheKey, &accessToken)
	if err == nil {
		return
	}

	// cache失效，从微信服务器获取
	var resAccessToken ResAccessToken
	var url = fmt.Sprintf(workAccessTokenURL, ak.handler.appID, ak.handler.appSecret)
	resAccessToken, err = GetTokenFromServer(url)
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 1500
	err = ak.handler.cache.SetTTL(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	accessToken = resAccessToken.AccessToken
	return
}

// GetTokenFromServer 强制从微信服务器获取token
func GetTokenFromServer(url string) (resAccessToken ResAccessToken, err error) {
	req := rest.New()
	code, err := req.GetJson(url, nil, &resAccessToken)
	if err != nil {
		return
	}

	if code != iris.StatusOK {
		err = fmt.Errorf("http code: %d", code)
		return
	}

	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}
	return
}
