package router

import (
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

func Router(root server.Router) {

	log.Info("init router...")
	api := root.Path("/api")
	api.Get("/", apiListSample)
	api.Get("/:id", apiOneAndIncreaseSample)
	api.Get("/:id/fail", apiTransactionFailSample)
}
