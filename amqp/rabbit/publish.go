package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
)

func publish(session *sessionRabbit, ch *amqp.Channel, m Message, options ...PublishOption) (err error) {

	//log.Debug("handle publish message:", m.Router.Exchange, m.Router.Key)
	if !session.isReady {
		return ErrNotConnected
	}
	if len(m.Router.Exchange) <= 0 {
		m.Router.Exchange = defaultExchange
	}

	var p amqp.Publishing
	if len(options) == 0 {
		options = append(options, defaultPublishOption)
	}

	for _, option := range options {
		p, err = option(p, m.Payload)
		if err != nil {
			return
		}
	}

	if err != nil {
		return
	}

	err = ch.Publish(m.Router.Exchange, m.Router.Key, false, false, p)
	if err != nil {
		log.Info("submit failed")
		log.Info(err)
		return
	}
	//log.Debug("submit completed")

	return err
}
