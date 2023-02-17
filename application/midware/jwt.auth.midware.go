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
	h, ok := GetAuth(ctx)
	if !ok {
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
		return
	}

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

func GetAuth(ctx iris.Context) (auth string, ok bool) {

	header := ctx.GetHeader(authHeader)
	if len(header) <= 0 {
		log.Debug("no auth")
		ok = false
		return
	}
	if !strings.HasPrefix(header, Realm) {
		log.Debug("unsupported realm:%s", header)
		ok = false
		return
	}
	auth = strings.ReplaceAll(header, Realm, "")
	return
}