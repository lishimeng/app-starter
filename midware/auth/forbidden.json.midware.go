package auth

import "github.com/kataras/iris/v12"

type ForbiddenType string

const (
	JsonForbiddenResponse ForbiddenType = "JSON"
	HtmlForbiddenResponse               = "HTML"
)

type ForbiddenOption struct {
	Type            ForbiddenType
	Redirect401Page string
}

type ForbiddenFunc func() ForbiddenOption

var WithJsonResp = func() ForbiddenOption {
	return ForbiddenOption{
		Type: JsonForbiddenResponse,
	}
}

var WithHtmlResp = func(redirect string) ForbiddenOption {
	return ForbiddenOption{
		Type:            HtmlForbiddenResponse,
		Redirect401Page: redirect,
	}
}

// Forbidden401Handler 401处理器
//
// 判定为401时返回消息
func Forbidden401Handler(option ...ForbiddenFunc) func(iris.Context) {
	var opt ForbiddenOption
	for _, f := range option {
		opt = f()
	}
	return func(ctx iris.Context) {
		if ctx.Values().Get(UserInfoKey) == nil {
			switch opt.Type {
			case HtmlForbiddenResponse: // 返回401页面
				ctx.Redirect(opt.Redirect401Page)
			case JsonForbiddenResponse: // 返回json
				errorWith(ctx, iris.StatusUnauthorized, ErrNotAllowed)
			default: // 返回http 401状态值
				ctx.StatusCode(iris.StatusUnauthorized)
			}
			return
		} else {
			ctx.Next()
		}
	}
}
