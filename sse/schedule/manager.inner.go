package schedule

import (
	"github.com/lishimeng/app-starter/sse/client"
	"github.com/lishimeng/app-starter/sse/event"
	"github.com/lishimeng/go-log"
)

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

// sendToClient 发送消息到指定客户端，payload中的event决定发送目标
// 如果客户端没有订阅该event，则忽略发送
func (m *Manager) sendToClient(clientId string, payload *event.Payload) error {
	message := payload.Marshall()
	if message == "" {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.Clients[clientId]
	if !ok {
		return nil
	}

	// 检查客户端是否订阅了payload中的event
	if !c.IsSubscribed(payload.Event) {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
		}
	}()

	err := c.SendMessage(message)
	if err != nil {
		log.Info(err)
	}
	return nil
}

// broadcast 广播消息给所有客户端，payload中的event决定发送目标
// 只有订阅了该event的客户端才会收到消息
func (m *Manager) broadcast(payload *event.Payload) error {
	message := payload.Marshall()
	if message == "" {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, c := range m.Clients {
		// 检查客户端是否订阅了payload中的event
		if !c.IsSubscribed(payload.Event) {
			continue
		}

		func(client *client.Client) {
			defer func() {
				if err := recover(); err != nil {
					log.Info(err)
				}
			}()
			err := client.SendMessage(message)
			if err != nil {
				log.Info(err)
			}
		}(c)
	}
	return nil
}
