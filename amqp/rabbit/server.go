package rabbit

import (
	"github.com/lishimeng/go-log"
	"time"
)

func (session *sessionRabbit) Publish(m Message) (err error) {
	defer func() {
		_ = recover()
	}()
	select {
	case session.globalTxChannel <- m:
		log.Debug("message to tx task")
		return
	case <-time.After(time.Millisecond * 200):
		err = ErrPublishTimeout
		return
	}

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
