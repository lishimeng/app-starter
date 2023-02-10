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
				log.Debug("wait conn ready [%ds]", t)
			})
		default:
			delay.Reset()
			session.globalChannelLoop()
		}
	}
}

// globalChannelLoop 不会出现panic
func (session *sessionRabbit) globalChannelLoop() {

	log.Debug("global tx loop start")
	defer func() {
		log.Debug("global tx loop exit")
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
			log.Debug("session ctx close")
			return
		case <-session.connCtx.Done():
			return
		case m, ok := <-session.globalTxChannel:
			if !ok {
				log.Debug("global tx has closed")
				return
			}
			// TODO 重发
			e := publish(session, ch, m, onPublished)
			if e != nil {
				log.Info(e) // 发送失败
			}
		}
	}
}
