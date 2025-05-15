package ws

import (
	"github.com/gorilla/websocket"
	"github.com/lishimeng/app-starter/server"
	"net/http"
)

type Conn struct {
	C *websocket.Conn
}

type Logic func(data any, txHandler TxHandler)

var Build = func(topic string, logic Logic) server.Handler {
	return func(ctx server.Context) {
		handleWsSession(ctx, topic, logic)
	}
}

var (
	upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
}
