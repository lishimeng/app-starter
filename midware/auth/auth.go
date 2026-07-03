package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/token"
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
	uid = ctx.GetHeader(UidKey)
	return
}

func GetClientType(ctx server.Context) (ct string) {
	ct = ctx.GetHeader(ClientKey)
	return
}

func GetOrg(ctx server.Context) (org string) {
	org = ctx.GetHeader(OrgKey)
	return
}

func GetDept(ctx server.Context) (dept string) {
	dept = ctx.GetHeader(DeptKey)
	return
}

func GetScope(ctx server.Context) (scope string) {
	scope = ctx.GetHeader(ScopeKey)
	return
}

func HasScope(ctx server.Context, expect string) (ok bool) {
	scope := ctx.GetHeader(ScopeKey)
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
	ui, ok := ctx.Get(UserInfoKey)
	if !ok {
		err = ErrNotAllowed
		return
	}
	uid, ok = ui.(token.JwtPayload)
	if !ok {
		err = ErrNotAllowed
		return
	}
	return
}

func unauthorized(ctx server.Context) {
	errorWith(ctx, http.StatusUnauthorized, ErrNotAllowed)
}
