package auth

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

// LoginRedirect 如果需要登录,跳转到登录界面
func LoginRedirect(redirect string) func(server.Context) {
	return func(ctx server.Context) {
		if _, ok := ctx.Get(UserInfoKey); !ok {
			ctx.Redirect(redirect, http.StatusFound)
			return
		}
		ctx.Next()
	}
}
