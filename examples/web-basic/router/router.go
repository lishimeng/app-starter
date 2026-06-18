package router

import (
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
)

func Router(root server.Router) {

	log.Info("init router...")
	root.Get("/", apiListSample)
	root.Get("/{id}", apiOneAndIncreaseSample)
	root.Get("/{id}/fail", apiTransactionFailSample)
}
