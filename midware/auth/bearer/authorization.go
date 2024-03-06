package bearer

import (
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
	"strings"
)

const (
	AuthHeader = "Authorization"
	Realm      = "Bearer "
)

func GetAuth(ctx server.Context) (auth string, ok bool) {

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
