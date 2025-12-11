package cookie

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
)

var HttpsCookie = false
var AccessTokenName = "at"
var RefreshTokenName = "rt"

func del(ctx server.Context, name string) {
	ctx.C.SetCookie(&iris.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   HttpsCookie,
	})
}

// RemoveAuth 清除token
func RemoveAuth(ctx server.Context) {
	del(ctx, AccessTokenName)
	del(ctx, RefreshTokenName)
}

func Save(ctx server.Context, name string, value string, ttlHour int, domain string) {
	sameSiteMode := http.SameSiteDefaultMode
	if HttpsCookie {
		sameSiteMode = http.SameSiteStrictMode
	}

	cookie := iris.Cookie{
		Name:     name,           // Cookie 名称
		Value:    value,          // 业务值（如 JWT）
		Path:     "/",            // 生效路径
		MaxAge:   3600 * ttlHour, // 有效期（秒），24小时
		HttpOnly: true,           // 核心：开启 HttpOnly
		Secure:   HttpsCookie,    // 本地测试设为 false，生产 HTTPS 必须设为 true
		SameSite: sameSiteMode,   // 防止 CSRF
		Domain:   domain,
	}
	ctx.C.SetCookie(&cookie)
}

// SaveAuth 保存token
func SaveAuth(ctx server.Context, token string, ttlHour int, domain string) {
	Save(ctx, AccessTokenName, token, ttlHour, domain)
}

// SaveRefreshAuth 保存refresh token
func SaveRefreshAuth(ctx server.Context, token string, ttlHour int, domain string) {
	Save(ctx, RefreshTokenName, token, ttlHour, domain)
}

// GetAuth 读token
func GetAuth(ctx server.Context) string {
	c := ctx.C.GetCookie(AccessTokenName)
	return c
}

// GetRefreshAuth 读refresh token
func GetRefreshAuth(ctx server.Context) string {
	c := ctx.C.GetCookie(RefreshTokenName)
	return c
}
