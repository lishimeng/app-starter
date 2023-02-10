package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
	"time"
)

func publish(session *sessionRabbit, ch *amqp.Channel, m Message, notifyPublish chan amqp.Confirmation) (err error) {

	//log.Debug("handle publish message:", m.Router.Exchange, m.Router.Key)

	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()

	if !session.isReady {
		return ErrNotConnected
	}
	if len(m.Router.Exchange) <= 0 {
		m.Router.Exchange = defaultExchange
	}

	var p amqp.Publishing
	if len(m.Options) == 0 {
		m.Options = append(m.Options, defaultPublishOption)
	}

	for _, option := range m.Options {
		p, err = option(p, m.Payload)
		if err != nil {
			return
		}
	}

	if err != nil {
		return
	}

	p.Timestamp = time.Now()

	//log.Debug("submit completed")

	for {
		err = ch.Publish(m.Router.Exchange, m.Router.Key, false, false, p)
		if err != nil {
			log.Info("submit failed")
			log.Info(err)
			return
		}

		select {
		case <-session.ctx.Done():
			return
		case <-session.connCtx.Done():
			return
		case confirm := <-notifyPublish:
			if confirm.Ack {
				return
			}
		case <-time.After(time.Second):

		}
	}

	return err
}
