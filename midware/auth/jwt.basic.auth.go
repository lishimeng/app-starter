package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
)

// JwtBasic 预处理jwt,解析后存入 UserInfoKey 和相应header
//
// 需要启动token验证器
func JwtBasic() func(iris.Context) {
	return func(ctx iris.Context) {
		var err error
		h, ok := tool.GetAuth(ctx)
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
			log.Debug("can't verify token")
			log.Debug(err)
			ctx.Next()
			return
		}

		ctx.Values().Set(UserInfoKey, p)

		r := ctx.Request()

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
		ctx.Next()
	}
}
