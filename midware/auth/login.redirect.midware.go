package auth

import "github.com/kataras/iris/v12"

// LoginRedirect 如果需要登录,跳转到登录界面
//
// 检查 UserInfoKey 是否存在,如果不存在,判定为未登录
//
// redirect: 登录地址,需要全路径
//
// 需要启动token验证器
func LoginRedirect(redirect string) func(iris.Context) {
	return func(ctx iris.Context) {
		ui := ctx.Values().Get(UserInfoKey)
		if ui == nil {
			ctx.Redirect(redirect, 302)
			return
		} else {
			ctx.Next()
		}
	}
}
