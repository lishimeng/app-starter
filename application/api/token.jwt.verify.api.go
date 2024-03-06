package api

import (
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/midware/auth/bearer"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/go-log"
)

func JwtTokenVerify(ctx server.Context) {

	var err error
	var resp token.HttpTokenResp

	authCtx, ok := bearer.GetAuth(ctx)
	if !ok {
		resp.Valid = false
		ctx.Json(resp)
		return
	}
	uid := ""
	err = factory.GetCache().Get(authCtx, &uid)
	if err != nil {
		log.Info(err)
		resp.Valid = false
		ctx.Json(resp)
		return
	}

	log.Info("%s", uid)
	resp.Valid = true
	ctx.Json(resp)
}
