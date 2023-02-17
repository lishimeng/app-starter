package midware

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/go-log"
	"strings"
)

const (
	authHeader  = "Authorization"
	Realm       = "Bearer "
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

	h := ctx.GetHeader(authHeader)
	if len(h) <= 0 {
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
		return
	}

	if !strings.HasPrefix(h, Realm) {
		log.Debug("unsupported realm:%s", h)
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
		return
	}

	h = strings.ReplaceAll(h, Realm, "")

	if TokenStorage == nil {
		log.Debug("token storage nil")
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
		return
	}

	p, err := TokenStorage.Verify(h)
	if err != nil {
		log.Debug("can't verify token")
		log.Debug(err)
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
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
