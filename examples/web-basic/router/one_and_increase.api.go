package router

import (
	app "github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/model"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
)

func apiOneAndIncreaseSample(ctx server.Context) {
	var err error
	id, err := ctx.C.Params().GetInt("id")
	if err != nil || id <= 0 {
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: 400, Message: "invalid id"},
		})
		return
	}

	var row model.BusinessConnector
	txErr := app.Transaction(func(tx persistence.TxContext) (e error) {
		if e = tx.Model(&model.BusinessConnector{}).Equal("id", id).First(&row); e != nil {
			return
		}
		row.Status++
		if e = tx.Model(&row).Update("status", row.Status); e != nil {
			return
		}
		return
	})
	if txErr != nil {
		code, msg := 500, txErr.Error()
		if persistence.IsNotFound(txErr) {
			code, msg = 404, "not found"
		}
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: code, Message: msg},
		})
		return
	}

	var dto businessConnectorDto
	toBusinessConnectorDto(row, &dto)
	ctx.Json(app.ResponseWrapper{
		Response: app.Response{Code: 0},
		Data:     dto,
	})
}
