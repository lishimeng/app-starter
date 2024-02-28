package token

import (
	"github.com/kataras/jwt"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/rest"
)

// httpStorage http token server
//
// 通过http接口验证token
type httpStorage struct {
	connector HttpStorageConnector
}

type HttpStorageConnector struct {
	Server string // 全路径 包含schema
}

const (
	// defaultTokenServer 默认token server
	defaultTokenServer = "http://token.default.svc.cluster.local/api/token/verify"
)

func NewHttpStorage(c HttpStorageConnector) (s Storage) {
	server := &httpStorage{connector: c}
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
	ut, err := jwt.Decode([]byte(key))
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
	var resp HttpTokenResp
	header, value := tool.BuildAuth(key)
	code, err := rest.New().GetJson(hs.connector.Server, nil, &resp, rest.Header{Name: header, Value: value})
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
