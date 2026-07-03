package auth

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

type ForbiddenType string

const (
	JsonForbiddenResponse ForbiddenType = "JSON"
	HtmlForbiddenResponse               = "HTML"
)

type ForbiddenOption struct {
	Type            ForbiddenType
	Redirect401Page string
	Scope           string
}

var WithJsonResp = func() func(ForbiddenOption) ForbiddenOption {
	return func(opt ForbiddenOption) ForbiddenOption {
		opt.Type = JsonForbiddenResponse
		return opt
	}
}

// WithScope 设置本程序需要的scope, 一个程序选择一个scope
var WithScope = func(scope string) func(ForbiddenOption) ForbiddenOption {
	return func(opt ForbiddenOption) ForbiddenOption {
		opt.Scope = scope
		return opt
	}
}

var WithHtmlResp = func(redirect string) func(ForbiddenOption) ForbiddenOption {
	return func(opt ForbiddenOption) ForbiddenOption {
		opt.Type = HtmlForbiddenResponse
		opt.Redirect401Page = redirect
		return opt
	}
}

// Forbidden401Handler 401处理器
func Forbidden401Handler(option ...func(ForbiddenOption) ForbiddenOption) func(server.Context) {
	var opt ForbiddenOption
	for _, f := range option {
		opt = f(opt)
	}
	return func(ctx server.Context) {
		if !checkForbidden(ctx, opt) {
			responseForbidden(ctx, opt)
			return
		}
		ctx.Next()
	}
}

func responseForbidden(ctx server.Context, opt ForbiddenOption) {
	switch opt.Type {
	case HtmlForbiddenResponse:
		ctx.Redirect(opt.Redirect401Page, http.StatusFound)
	case JsonForbiddenResponse:
		unauthorized(ctx)
	default:
		ctx.Status(http.StatusUnauthorized)
	}
}

func checkForbidden(ctx server.Context, opt ForbiddenOption) (pass bool) {
	pass = true
	if _, ok := ctx.Get(UserInfoKey); !ok {
		pass = false
		return
	}
	grantedScope := ctx.GetHeader(ScopeKey)
	if len(opt.Scope) > 0 {
		pass = checkScope(grantedScope, opt.Scope)
	}
	return
}

func checkScope(scope string, expected string) bool {
	return hasScope(scope, expected)
}
