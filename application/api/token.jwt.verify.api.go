package api

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
)

func JwtTokenVerify(ctx iris.Context) {

	var err error
	var resp token.HttpTokenResp

	auth, ok := tool.GetAuth(ctx)
	if !ok {
		resp.Valid = false
		tool.ResponseJSON(ctx, resp)
		return
	}
	uid := ""
	err = factory.GetCache().Get(auth, &uid)
	if err != nil {
		log.Info(err)
		resp.Valid = false
		tool.ResponseJSON(ctx, resp)
		return
	}

	log.Info("%s", uid)
	resp.Valid = true
	tool.ResponseJSON(ctx, resp)
}
