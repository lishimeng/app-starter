package auth

import (
	"github.com/lishimeng/app-starter/midware/auth/bearer"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

// SimpleJwtAuth 简单验证器
//
// 判定无权限后返回json类型message和http401.如有权限,将数据存入 UserInfoKey 和相应header
//
// 需要启动token验证器
func SimpleJwtAuth(ctx server.Context) {
	h, ok := bearer.GetAuth(ctx)
	if !ok {
		unauthorized(ctx)
		return
	}

	if TokenStorage == nil {
		log.Debug("token storage nil")
		unauthorized(ctx)
		return
	}

	p, err := TokenStorage.Verify(h)
	if err != nil {
		log.Debug("can't verify token", "err", err)
		unauthorized(ctx)
		return
	}

	ctx.Set(UserInfoKey, p)

	r := ctx.Request()
	if r == nil {
		unauthorized(ctx)
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
