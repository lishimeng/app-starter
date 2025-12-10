package bearer

import (
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
)

const (
	AuthHeader = "Authorization"
	Realm      = "Bearer "
)

var ForceCookie = false
var HttpsCookie = false
var CookieName = "token"

func getCookieAuth(ctx server.Context) string {
	c := ctx.C.GetCookie(CookieName)
	return c
}

func RemoveCookieAuth(ctx server.Context) {
	ctx.C.SetCookie(&iris.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   HttpsCookie,
	})
}

func SaveCookieAuth(ctx server.Context, token string, ttlHour int, domain string) {
	sameSiteMode := http.SameSiteDefaultMode
	if HttpsCookie {
		sameSiteMode = http.SameSiteStrictMode
	}
	if len(domain) > 0 && !strings.HasPrefix(domain, ".") {
		domain = "." + domain
	}
	cookie := iris.Cookie{
		Name:     CookieName,     // Cookie 名称
		Value:    token,          // 业务值（如 JWT）
		Path:     "/",            // 生效路径
		MaxAge:   3600 * ttlHour, // 有效期（秒），24小时
		HttpOnly: true,           // 核心：开启 HttpOnly
		Secure:   HttpsCookie,    // 本地测试设为 false，生产 HTTPS 必须设为 true
		SameSite: sameSiteMode,   // 防止 CSRF
		Domain:   domain,
	}
	ctx.C.SetCookie(&cookie)
}

func GetAuth(ctx server.Context) (auth string, ok bool) {

	auth = getCookieAuth(ctx)
	if len(auth) > 0 {
		ok = true
		return
	}
	if ForceCookie {
		auth = ""
		ok = false
		return
	}
	header := ctx.C.GetHeader(AuthHeader)
	if len(header) <= 0 {
		log.Debug("no auth")
		ok = false
		return
	}
	if !strings.HasPrefix(header, Realm) {
		log.Debug("unsupported realm:%s", header)
		ok = false
		return
	}
	ok = true
	auth = strings.ReplaceAll(header, Realm, "")
	return
}

func BuildAuth(jwtStr string) (header string, value string) {
	header = AuthHeader
	value = Realm + jwtStr
	return
}
