package auth

import (
	"github.com/lishimeng/app-starter/server"
)

// LoginRedirect 如果需要登录,跳转到登录界面
//
// 检查 UserInfoKey 是否存在,如果不存在,判定为未登录
//
// redirect: 登录地址,需要全路径
//
// 需要启动token验证器
func LoginRedirect(redirect string) func(server.Context) {
	return func(ctx server.Context) {
		ui := ctx.C.Values().Get(UserInfoKey)
		if ui == nil {
			ctx.C.Redirect(redirect, 302)
			return
		} else {
			ctx.C.Next()
		}
	}
}
