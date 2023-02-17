package tool

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/go-log"
	"strings"
)

const (
	AuthHeader = "Authorization"
	Realm      = "Bearer "
)

func GetAuth(ctx iris.Context) (auth string, ok bool) {

	header := ctx.GetHeader(AuthHeader)
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
	auth = strings.ReplaceAll(header, Realm, "")
	return
}

func BuildAuth(jwtStr string) (header string, value string) {
	header = AuthHeader
	value = Realm + jwtStr
	return
}
