package rabbit

import "github.com/lishimeng/go-log"

func (session *sessionRabbit) Publish(m Message) {
	if !session.isReady {
		log.Debug("session is unready")
		return
	}
	session.globalChannel <- m
}

// Close will cleanly shut down the channel and conn.
func (session *sessionRabbit) Close() error {
	var err = session.disposeResources()
	session.isReady = false
	return err
}

func (session *sessionRabbit) Subscribe(r Route, rxHandler RxHandler) {
	subscribe(session, r, rxHandler)
}
