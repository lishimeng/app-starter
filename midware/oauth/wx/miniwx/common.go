package miniwx

const (
	// AccessTokenURL 获取access_token的接口
	credentialsURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// AccessTokenURL 获取access_token的接口
	loginURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// ResAccessToken struct
type ResAccessToken struct {
	CommonError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type WxMiniLoginResp struct {
	SessionKey string `json:"session_Key,omitempty"`
	UnionId    string `json:"unionid,omitempty"`
	OpenId     string `json:"openid,omitempty"`
	CommonError
}
