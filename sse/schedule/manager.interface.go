package schedule

import (
	"context"

	"github.com/lishimeng/app-starter/sse/event"
)

func (m *Manager) Ctx() context.Context {
	return m.ctx
}

func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.Clients)
}

// GetClientsByEvent 获取订阅了指定event的所有客户端ID
func (m *Manager) GetClientsByEvent(eventType string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []string
	for id, c := range m.Clients {
		if c.IsSubscribed(eventType) {
			result = append(result, id)
		}
	}
	return result
}

// GetClientEvents 获取指定客户端订阅的所有event
func (m *Manager) GetClientEvents(clientId string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if c, ok := m.Clients[clientId]; ok {
		return c.GetEvents()
	}
	return nil
}

// SendEvent 发送事件
func (m *Manager) SendEvent(e event.Event) (err error) {
	if e.Payload == nil {
		return
	}
	switch e.Type {
	case event.ToClient:
		if len(e.ClientId) == 0 {
			return
		}
		err = m.sendToClient(e.ClientId, e.Payload)
	case event.Broadcast:
		err = m.broadcast(e.Payload)
	}
	return
}
