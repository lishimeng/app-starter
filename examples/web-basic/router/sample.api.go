package router

import (
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/model"
	"github.com/lishimeng/app-starter/server"
)

func apiSample(ctx server.Context) {
	var resp = make(map[string]any)
	var err error

	var list []model.BusinessConnector
	_, err = app.GetOrm().Query(new(model.BusinessConnector)).Limit(10).All(&list)
	if err != nil {
		resp["err"] = err.Error()
		ctx.Json(resp)
		return
	}
	resp["a"] = "b"
	resp["list"] = list
	ctx.Json(resp)
}
