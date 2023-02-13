package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
	"time"
)

func publish(session *sessionRabbit, ch *amqp.Channel, m Message, notifyPublish chan amqp.Confirmation) (err error) {

	log.Fine("handle publish message to exchange:[%s]key:[%s]", m.Router.Exchange, m.Router.Key)

	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()

	if !session.isReady {
		log.Debug("session is unready, can't publish the message")
		return ErrNotConnected
	}
	if len(m.Router.Exchange) <= 0 {
		log.Debug("use default exchange:%s", defaultExchange)
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
	p.DeliveryMode = amqp.Persistent

	err = doPublish(session, ch, m, notifyPublish, p)
	if err != nil {
		log.Info(err)
	}

	return
}

func doPublish(session *sessionRabbit, ch *amqp.Channel, m Message, notifyPublish chan amqp.Confirmation, p amqp.Publishing) (err error) {
	log.Fine("publish message to server:%s[%s]", p.MessageId, p.Timestamp.String())
	var publishTimes = 0
	for {
		publishTimes++
		if publishTimes > 1 {
			log.Debug("publish times:%d, %s[%s]", publishTimes, p.MessageId, p.Timestamp.String())
		}

		err = ch.Publish(m.Router.Exchange, m.Router.Key, false, false, p)
		if err != nil {
			log.Info("publish failed")
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
				log.Fine("server confirmed:%s", p.MessageId)
				return
			} else {
				log.Debug("server not receive:%s", p.MessageId)
			}
		case <-time.After(time.Millisecond * 1):
			log.Info("publish confirm timeout:%s", p.MessageId)
		}
	}
}
