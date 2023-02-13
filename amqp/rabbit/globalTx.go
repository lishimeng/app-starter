package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
)

func (session *sessionRabbit) globalChannelProcess() {

	var delay = NewDelay(1, 60, true)
	for {
		select {
		case <-session.ctx.Done():
			return
		case <-session.connCtx.Done():
			delay.Delay(func(t int) {
				log.Fine("wait conn ready [%ds]", t)
			})
		default:
			delay.Reset()
			session.globalChannelLoop()
		}
	}
}

// globalChannelLoop 不会出现panic
func (session *sessionRabbit) globalChannelLoop() {

	log.Fine("global tx loop start")
	defer func() {
		log.Fine("global tx loop exit")
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
			log.Fine("session ctx close")
			return
		case <-session.connCtx.Done():
			return
		case m, ok := <-session.globalTxChannel:
			if !ok {
				log.Debug("global tx has closed")
				return
			}
			e := publish(session, ch, m, onPublished)
			if e != nil {
				log.Info(e) // 发送失败
			}
		}
	}
}
