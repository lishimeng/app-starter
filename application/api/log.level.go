package api

import (
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

type LogLevelReq struct {
	Level string `json:"level,omitempty"`
}

type Resp struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func changeLogLevel(ctx server.Context) {
	var err error
	var req LogLevelReq
	var resp Resp
	err = ctx.C.ReadJSON(&req)
	if err != nil {
		resp.Code = 500
		resp.Message = err.Error()
		ctx.Json(resp)
		return
	}
	if req.Level == "" {
		resp.Code = 200
		resp.Message = "unknown level"
		ctx.Json(resp)
		return
	}
	if err = log.SetLevelFromString(req.Level); err != nil {
		resp.Code = 200
		resp.Message = "unknown level"
		ctx.Json(resp)
		return
	}

	resp.Code = 200
	ctx.Json(resp)
}
