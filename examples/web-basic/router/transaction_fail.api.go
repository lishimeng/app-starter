package router

import (
	"errors"
	"fmt"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/model"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/server"
	"gorm.io/gorm"
)

type transactionFailDto struct {
	Id              int    `json:"id"`
	StatusBefore    int    `json:"statusBefore"`
	AttemptedStatus int    `json:"attemptedStatus"`
	StatusAfter     int    `json:"statusAfter"`
	RolledBack      bool   `json:"rolledBack"`
	TxError         string `json:"txError"`
}

func apiTransactionFailSample(ctx server.Context) {
	var err error
	id, err := ctx.C.Params().GetInt("id")
	if err != nil || id <= 0 {
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: 400, Message: "invalid id"},
		})
		return
	}

	var before model.BusinessConnector
	if err = app.Query(func(o persistence.OrmContext) error {
		return o.Model(&model.BusinessConnector{}).Equal("id", id).First(&before)
	}); err != nil {
		code, msg := 500, err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code, msg = 404, "not found"
		}
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: code, Message: msg},
		})
		return
	}

	var attemptedStatus int
	txErr := app.Transaction(func(tx persistence.TxContext) (e error) {
		var row model.BusinessConnector
		if e = tx.Model(&model.BusinessConnector{}).Equal("id", id).First(&row); e != nil {
			return
		}
		row.Status++
		attemptedStatus = row.Status
		if e = tx.Model(&row).Update("status", row.Status); e != nil {
			return
		}
		e = fmt.Errorf("deliberate rollback for demo")
		return
	})
	if txErr == nil {
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: 500, Message: "expected transaction rollback, but committed"},
		})
		return
	}
	log.With("err", txErr).Error("transaction fail")
	log.Errorf("%v", txErr)

	var after model.BusinessConnector
	if err = app.Query(func(o persistence.OrmContext) error {
		return o.Model(&model.BusinessConnector{}).Equal("id", id).First(&after)
	}); err != nil {
		ctx.Json(app.ResponseWrapper{
			Response: app.Response{Code: 500, Message: err.Error()},
		})
		return
	}

	rolledBack := after.Status == before.Status && attemptedStatus == before.Status+1
	ctx.Json(app.ResponseWrapper{
		Response: app.Response{Code: 0, Message: "transaction failed and rolled back"},
		Data: transactionFailDto{
			Id:              id,
			StatusBefore:    before.Status,
			AttemptedStatus: attemptedStatus,
			StatusAfter:     after.Status,
			RolledBack:      rolledBack,
			TxError:         txErr.Error(),
		},
	})
}
