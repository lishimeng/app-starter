package auth

import (
	"github.com/kataras/iris/v12"
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
//
// 判定为401时返回消息
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
		ctx.C.Next()
	}
}

func responseForbidden(ctx server.Context, opt ForbiddenOption) {
	switch opt.Type {
	case HtmlForbiddenResponse: // 返回401页面
		ctx.C.Redirect(opt.Redirect401Page)
	case JsonForbiddenResponse: // 返回json
		errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
	default: // 返回http 401状态值
		ctx.C.StatusCode(iris.StatusUnauthorized)
	}
}

func checkForbidden(ctx server.Context, opt ForbiddenOption) (pass bool) {
	pass = true
	// 检查不通过的情况
	// ui为空
	if ctx.C.Values().Get(UserInfoKey) == nil {
		pass = false
		return
	}
	// scope检查
	grantedScope := ctx.C.GetHeader(ScopeKey)
	if len(opt.Scope) > 0 {
		pass = checkScope(grantedScope, opt.Scope)
	}

	return
}

func checkScope(scope string, expected string) bool {
	return hasScope(scope, expected)
}
