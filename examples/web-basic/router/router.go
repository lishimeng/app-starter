package router

import "github.com/lishimeng/app-starter/server"

func Router(root server.Router) {

	root.Get("/", apiListSample)
	root.Get("/{id}", apiOneAndIncreaseSample)
}
