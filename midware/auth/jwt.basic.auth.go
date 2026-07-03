package auth

import (
	"github.com/lishimeng/app-starter/midware/auth/bearer"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

// JwtBasic 预处理jwt,解析后存入 UserInfoKey 和相应header
//
// 需要启动token验证器
func JwtBasic() func(server.Context) {
	return func(ctx server.Context) {
		h, ok := bearer.GetAuth(ctx)
		if !ok {
			ctx.Next()
			return
		}

		if TokenStorage == nil {
			log.Debug("token storage nil")
			ctx.Next()
			return
		}

		p, err := TokenStorage.Verify(h)
		if err != nil {
			log.Debug("can't verify token", "err", err)
			ctx.Next()
			return
		}

		ctx.Set(UserInfoKey, p)

		r := ctx.Request()
		if r == nil {
			ctx.Next()
			return
		}

		if len(p.Uid) > 0 {
			r.Header.Set(UidKey, p.Uid)
		}
		if len(p.Client) > 0 {
			r.Header.Set(ClientKey, p.Client)
		}
		if len(p.Org) > 0 {
			r.Header.Set(OrgKey, p.Org)
		}
		if len(p.Dept) > 0 {
			r.Header.Set(DeptKey, p.Dept)
		}
		if len(p.Scope) > 0 {
			r.Header.Set(ScopeKey, p.Scope)
		}
		ctx.Next()
	}
}
