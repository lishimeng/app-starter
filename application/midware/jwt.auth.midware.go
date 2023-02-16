package midware

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/go-log"
)

const (
	jwtHeader = "Auth"

	OrgKey      = "org"
	DeptKey     = "dept"
	ClientKey   = "clientType"
	UidKey      = "uid"
	UserInfoKey = "ui"
)

var (
	ErrNotAllowed = errors.New("401 not allowed")
)

var TokenStorage token.Storage

// JwtAuth 验证器， EnableTokenValidator后可用
func JwtAuth(ctx iris.Context) {

	var err error
	h := ctx.GetHeader(jwtHeader)
	if len(h) <= 0 {
		errorWith(ctx, 401, ErrNotAllowed)
		return
	}
	if TokenStorage == nil {
		log.Debug("token storage nil")
		errorWith(ctx, 401, ErrNotAllowed)
		return
	}

	p, err := TokenStorage.Verify(h)
	if err != nil {
		log.Debug("can't verify token")
		log.Debug(err)
		errorWith(ctx, 401, ErrNotAllowed)
		return
	}

	ctx.Values().Set(UserInfoKey, p)

	if len(p.Uid) > 0 {
		ctx.Header(UidKey, p.Uid)
	}
	if len(p.Client) > 0 {
		ctx.Header(ClientKey, p.Client)
	}
	if len(p.Org) > 0 {
		ctx.Header(OrgKey, p.Org)
	}
	if len(p.Dept) > 0 {
		ctx.Values().Set(DeptKey, p.Dept)
	}
	ctx.Next()
}