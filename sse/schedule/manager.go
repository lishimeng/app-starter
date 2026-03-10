package schedule

import (
	"context"
	"sync"

	"github.com/lishimeng/app-starter/sse/client"
	"github.com/lishimeng/go-log"
)

type Manager struct {
	ctx        context.Context
	mu         sync.RWMutex
	running    bool
	Clients    map[string]*client.Client // 客户端集合
	Register   chan *client.Client       // 注册客户端通道
	Unregister chan *client.Client       // 注销客户端通道
	Broadcast  chan string               // 广播消息通道
}

func New(ctx context.Context) *Manager {
	var m = &Manager{
		ctx:        ctx,
		Clients:    make(map[string]*client.Client),
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
		Broadcast:  make(chan string),
	}
	return m
}

func (m *Manager) Ctx() context.Context {
	return m.ctx
}

func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.Clients)
}

func (m *Manager) registerClient(c *client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Clients[c.ID] = c
	log.Info("客户端 %s 已连接，当前在线数: %d", c.ID, len(m.Clients))
}
func (m *Manager) unregisterClient(c *client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Clients[c.ID]; ok {
		delete(m.Clients, c.ID)
		log.Info("客户端 %s 已断开，当前在线数: %d", c.ID, len(m.Clients))
	}
}

func (m *Manager) broadcast(message string) {
	log.Info("broadcast: %s", message)
	for _, c := range m.Clients {
		// 按SSE协议格式写入消息
		m.sendTo(c, message)
	}
}

func (m *Manager) SendTo(clientId string, message string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
		}
	}()
	if c, ok := m.Clients[clientId]; ok {
		err := c.SendMessage(message)
		if err != nil {
			log.Info(err)
		}
	}
}

func (m *Manager) sendTo(c *client.Client, message string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
		}
	}()
	err := c.SendMessage(message)
	if err != nil {
		log.Info(err)
	}
}
