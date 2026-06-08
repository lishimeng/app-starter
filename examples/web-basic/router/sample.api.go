package router

import "github.com/lishimeng/app-starter/server"

func apiSample(ctx server.Context) {
	var resp = make(map[string]any)
	resp["a"] = "b"
	ctx.Json(resp)
}
