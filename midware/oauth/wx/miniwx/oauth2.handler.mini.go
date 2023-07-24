package miniwx

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/rest"
	"time"
)

type MiniAccessToken struct {
	handler BaseHandler
}

func New(opts ...Option) Oauth2Handle {
	handler := newBaseHandler()
	for _, opt := range opts {
		opt(&handler)
	}
	return &MiniAccessToken{
		handler: handler,
	}
}

func (ak *MiniAccessToken) Credentials() (accessToken string, err error) {
	// 双检，防止重复从微信服务器获取
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.handler.cacheKeyPrefix, ak.handler.appID)
	if ak.handler.cacheEnable {
		if err = ak.handler.cache.Get(accessTokenCacheKey, &accessToken); err == nil {
			return
		}
	}

	// cache中没有，从微信服务器获取
	var resAccessToken ResAccessToken
	var url = fmt.Sprintf(credentialsURL, ak.handler.appID, ak.handler.appSecret)
	err = GetTokenFromServer(url, resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
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

// AuthorizeCode 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *MiniAccessToken) AuthorizeCode(code string) (accessToken WxMiniLoginResp, err error) {

	var url = fmt.Sprintf(loginURL, ak.handler.appID, ak.handler.appSecret, code)
	err = GetTokenFromServer(url, &accessToken)
	if err != nil {
		return
	}

	if accessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", accessToken.ErrCode, accessToken.ErrMsg)
		return
	}
	return
}

// GetTokenFromServer 强制从微信服务器获取token
func GetTokenFromServer(url string, resp interface{}) (err error) {
	req := rest.New()
	code, err := req.GetJson(url, nil, &resp)
	if err != nil {
		return
	}

	if code != iris.StatusOK {
		err = fmt.Errorf("http code: %d", code)
		return
	}

	return
}
