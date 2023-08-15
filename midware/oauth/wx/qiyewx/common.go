package qiyewx

const (
	// AccessTokenURL 获取access_token的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// AccessTokenURL 获取access_token的接口
	miniLoginURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	// AccessTokenURL 企业微信获取access_token的接口
	workAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// CacheKeyOfficialAccountPrefix 微信公众号cache key前缀
	CacheKeyOfficialAccountPrefix = "gowechat_officialaccount_"
	// CacheKeyMiniProgramPrefix 小程序cache key前缀
	CacheKeyMiniProgramPrefix = "gowechat_miniprogram_"
	// CacheKeyWorkPrefix 企业微信cache key前缀
	CacheKeyWorkPrefix = "gowechat_work_"
)

// ResAccessToken struct
type ResAccessToken struct {
	CommonError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type CommonError struct {
	apiName string
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
