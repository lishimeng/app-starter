package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
)

func (session *sessionRabbit) txProcess(name string) {

	var delay = NewDelay(1, 60, true)
	for {
		select {
		case <-session.ctx.Done():
			return
		case <-session.connCtx.Done():
			delay.Delay(func(t int) {
				log.Fine("wait conn ready [%s][%ds]", name, t)
			})
		default:
			delay.Reset()
			session.txLoop(name)
		}
	}
}

// txLoop 不会出现panic
func (session *sessionRabbit) txLoop(name string) {

	log.Fine("global tx loop start:%s", name)
	defer func() {
		log.Fine("global tx loop exit:%s", name)
	}()
	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()

	var ch, ok, err = session.openChannel()
	if err != nil {
		return
	}

	if ok {
		defer func() {
			_ = ch.Close()
		}()
	}

	err = ch.Confirm(false)
	if err != nil {
		return
	}

	var onPublished = make(chan amqp.Confirmation)
	ch.NotifyPublish(onPublished)

	for {
		select {
		case <-session.ctx.Done():
			log.Fine("session ctx close:%s", name)
			return
		case <-session.connCtx.Done():
			return
		case m, ok := <-session.globalTxChannel:
			if !ok {
				log.Debug("global tx has closed:%s", name)
				return
			}
			log.Debug("publish message:%s[%+v]", name, m.Router)
			e := publish(session, ch, m, onPublished)
			if e != nil {
				log.Info(e) // 发送失败
			}
		}
	}
}
