package sse

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/lishimeng/app-starter/sse/client"
	"github.com/lishimeng/app-starter/sse/schedule"
	"github.com/lishimeng/go-log"
)

func NewManager(ctx context.Context) *schedule.Manager {
	var m = schedule.New(ctx)
	return m
}

func WebHandler(manager *schedule.Manager, w http.ResponseWriter, r *http.Request, events ...string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Retry-After", "12000")

	clientID := strconv.FormatInt(time.Now().UnixNano(), 10)

	c := client.New(manager.Ctx(), clientID, r, w)

	if len(events) > 0 {
		c.Subscribe(events...)
		log.Info("客户端 %s 订阅events: %v", clientID, events)
	}

	manager.Register <- c
	defer func() {
		log.Info("客户端 %s handler退出，执行反注册", clientID)
		manager.Unregister <- c
	}()

	c.Run(time.Second * 30)
}
