package cookie

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

var HttpsCookie = false
var AccessTokenName = "at"
var RefreshTokenName = "rt"

func del(ctx server.Context, name string) {
	ctx.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   HttpsCookie,
	})
}

func RemoveAuth(ctx server.Context) {
	del(ctx, AccessTokenName)
	del(ctx, RefreshTokenName)
}

func Save(ctx server.Context, name string, value string, ttlHour int, domain string) {
	sameSiteMode := http.SameSiteDefaultMode
	if HttpsCookie {
		sameSiteMode = http.SameSiteStrictMode
	}

	ctx.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   3600 * ttlHour,
		HttpOnly: true,
		Secure:   HttpsCookie,
		SameSite: sameSiteMode,
		Domain:   domain,
	})
}

func SaveAuth(ctx server.Context, token string, ttlHour int, domain string) {
	Save(ctx, AccessTokenName, token, ttlHour, domain)
}

func SaveRefreshAuth(ctx server.Context, token string, ttlHour int, domain string) {
	Save(ctx, RefreshTokenName, token, ttlHour, domain)
}

func GetAuth(ctx server.Context) string {
	c, _ := ctx.GetCookie(AccessTokenName)
	return c
}

func GetRefreshAuth(ctx server.Context) string {
	c, _ := ctx.GetCookie(RefreshTokenName)
	return c
}
