package admin

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

const DemoPath = "/demo/ping"

// Register adds sample routes on the admin listener (:6060).
func Register(root server.Router) {
	root.Get(DemoPath, func(ctx server.Context) {
		ctx.JSON(map[string]any{
			"ok":      true,
			"service": "web-basic-admin",
		})
	})
	root.Get("/demo/health", func(ctx server.Context) {
		ctx.Status(http.StatusOK)
		_, _ = ctx.Write([]byte("admin-ok"))
	})
}
