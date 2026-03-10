package schedule

import (
	"context"
	"sync"

	"github.com/lishimeng/app-starter/sse/client"
)

type Manager struct {
	ctx        context.Context
	mu         sync.RWMutex
	running    bool
	Clients    map[string]*client.Client // 客户端集合
	Register   chan *client.Client       // 注册客户端通道
	Unregister chan *client.Client       // 注销客户端通道
}

func New(ctx context.Context) *Manager {
	var m = &Manager{
		ctx:        ctx,
		Clients:    make(map[string]*client.Client),
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
	}
	return m
}
