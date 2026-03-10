package schedule

func (m *Manager) Run() {
	if m.running {
		return
	}
	m.running = true
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				return
			case client := <-m.Register:
				m.registerClient(client)
			case client := <-m.Unregister:
				m.unregisterClient(client)
			case message := <-m.Broadcast:
				m.broadcast(message)
			}
		}
	}()
}
