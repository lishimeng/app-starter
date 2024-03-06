package api

import (
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
	"strings"
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
	var lvl log.Level
	if strings.HasPrefix(req.Level, "INFO") {
		lvl = log.INFO
	} else if strings.HasPrefix(req.Level, "FINE") {
		lvl = log.FINE
	} else if strings.HasPrefix(req.Level, "DEBUG") {
		lvl = log.DEBUG
	} else if strings.HasPrefix(req.Level, "ERROR") {
		lvl = log.ERROR
	} else {
		resp.Code = 200
		resp.Message = "unknown level"
		ctx.Json(resp)
		return
	}

	log.SetLevelAll(lvl)

	resp.Code = 200
	ctx.Json(resp)
}
