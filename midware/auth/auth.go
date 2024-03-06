package auth

import (
	"errors"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/token"
	"strings"
)

const (
	OrgKey      = "org"
	DeptKey     = "dept"
	ClientKey   = "clientType"
	UidKey      = "uid"
	UserInfoKey = "ui"
	ScopeKey    = "auth_scope"
)

var (
	ErrNotAllowed = errors.New("401 not allowed")
)

var TokenStorage token.Storage

type Response struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
}

func errorWith(ctx server.Context, code int, err error) {
	var resp Response
	resp.Code = code
	resp.Message = err.Error()
	ctx.Json(resp)
}

func GetUid(ctx server.Context) (uid string) {
	uid = ctx.C.GetHeader(UidKey)
	return
}

func GetClientType(ctx server.Context) (ct string) {
	ct = ctx.C.GetHeader(ClientKey)
	return
}

func GetOrg(ctx server.Context) (org string) {
	org = ctx.C.GetHeader(OrgKey)
	return
}

func GetDept(ctx server.Context) (dept string) {
	dept = ctx.C.GetHeader(DeptKey)
	return
}

func GetScope(ctx server.Context) (scope string) {
	scope = ctx.C.GetHeader(ScopeKey)
	return
}

// HasScope 检查是否包含了指定的scope
func HasScope(ctx server.Context, expect string) (ok bool) {
	scope := ctx.C.GetHeader(ScopeKey)
	return hasScope(scope, expect)
}

func hasScope(scopeHeader string, expect string) (ok bool) {
	scope := scopeHeader
	ss := strings.Split(scope, ",")
	for _, v := range ss {
		if expect == v {
			ok = true
			return
		}
	}
	return
}

func GetUserInfo(ctx server.Context) (uid token.JwtPayload, err error) {
	var ui = ctx.C.Values().Get(UserInfoKey)
	uid, ok := ui.(token.JwtPayload)
	if !ok {
		err = ErrNotAllowed
		return
	}
	return
}
