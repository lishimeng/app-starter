package router

import (
	"strconv"

	app "github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/model"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
)

type businessConnectorDto struct {
	Id       int     `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	ConnType string  `json:"connType"`
	Enabled  int     `json:"enabled"`
	Status   int     `json:"status"`
	Desc     *string `json:"desc,omitempty"`
}

func toBusinessConnectorDto(src model.BusinessConnector, dst *businessConnectorDto) {
	dst.Id = src.Id
	dst.Code = src.Code
	dst.Name = src.Name
	dst.ConnType = src.ConnType
	dst.Enabled = src.Enabled
	dst.Status = src.Status
	dst.Desc = src.Desc
}

func apiListSample(ctx server.Context) {
	pageNum, pageSize := 0, 0
	if v, err := ctx.C.URLParamInt("pageNum"); err == nil && v > 0 {
		pageNum = v
	}
	if v, err := ctx.C.URLParamInt("pageSize"); err == nil && v > 0 {
		pageSize = v
	}
	code := ctx.C.URLParam("code")
	name := ctx.C.URLParam("name")
	connType := ctx.C.URLParam("connType")
	enabledStr := ctx.C.URLParam("enabled")

	pager := &app.SimplePager[model.BusinessConnector, businessConnectorDto]{
		Pager: app.Pager[businessConnectorDto]{
			BasePager: app.BasePager{
				PageNum:  pageNum,
				PageSize: pageSize,
			},
		},
		OrderExp: []string{"id desc"},
		QueryBuilder: func(tx persistence.TxContext) persistence.Query {
			q := tx.Model(new(model.BusinessConnector)).
				EqualStr("code", code).
				ILikeStr("name", name).
				EqualStr("conn_type", connType)
			if enabledStr != "" {
				if enabled, err := strconv.Atoi(enabledStr); err == nil {
					q = q.Equal("enabled", enabled)
				}
			}
			return q
		},
		Transform: func(src model.BusinessConnector, dst *businessConnectorDto) {
			toBusinessConnectorDto(src, dst)
		},
	}

	if err := app.QueryPage(pager); err != nil {
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: 500, Message: err.Error()},
		})
		return
	}

	ctx.Json(app.ResponseWrapper{
		Response: app.Response{Code: 0},
		Data:     pager.Pager,
	})
}
