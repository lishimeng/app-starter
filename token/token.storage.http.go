package token

import (
	"github.com/lishimeng/app-starter/rest"
	"github.com/lishimeng/go-log"
)

// httpStorage http token server
//
// 通过http接口验证token
type httpStorage struct {
	connector   HttpStorageConnector
	jwtProvider *JwtProvider
}

type HttpStorageConnector struct {
	Server string // 全路径 包含schema
}

const (
	// defaultTokenServer 默认token server
	defaultTokenServer = "http://token.default.svc.cluster.local/api/token/verify"
)

func NewHttpStorage(c HttpStorageConnector, provider *JwtProvider) (s Storage) {
	server := &httpStorage{connector: c, jwtProvider: provider}
	if len(server.connector.Server) == 0 {
		server.connector.Server = defaultTokenServer
	}
	s = server
	return
}

type HttpTokenReq struct {
	Token string `json:"token,omitempty"`
}

type HttpTokenResp struct {
	Valid bool `json:"valid,omitempty"`
}

func (hs *httpStorage) Verify(key string) (p JwtPayload, err error) {
	ut, err := hs.jwtProvider.Decode([]byte(key))
	if err != nil {
		return
	}
	err = ut.Claims(&p)
	if err != nil {
		return
	}
	// verify from server
	err = hs.verify(key)
	return
}

func (hs *httpStorage) verify(key string) (err error) {
	var req = HttpTokenReq{Token: key}
	var resp HttpTokenResp
	code, err := rest.New().PostJson(hs.connector.Server, req, resp)
	if err != nil {
		log.Info(err)
		err = ErrInvalid
		return
	}
	if code != 200 {
		log.Debug("http status code:%d", code)
		err = ErrInvalid
		return
	}
	if !resp.Valid {
		log.Debug("token is not valid")
		err = ErrInvalid
		return
	}
	return
}
