package rabbit

import (
	"errors"
	"fmt"
	"github.com/lishimeng/go-log"
)

func (session *sessionRabbit) Publish(m Message) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
			log.Fine(err)
		}
	}()
	select {
	case session.globalTxChannel <- m:
		log.Fine("message to tx task")
		return
	default:
		err = ErrTxBufferFull
		return err
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
