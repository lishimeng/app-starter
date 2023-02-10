package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
)

func subscribe(session *sessionRabbit, r Route, rxHandler RxHandler) {

	if len(r.Exchange) <= 0 {
		r.Exchange = defaultExchange
	}

	var delay = NewDelay(1, 60, false)

	for {
		select {
		case <-session.ctx.Done(): // session销毁
			log.Info("subscribe exit [%s:%s-%s]", r.Exchange, r.Key, r.Queue)
			return
		case <-session.connCtx.Done(): // 连接断开
			// 等待连接
			delay.Delay(func(t int) {
				log.Debug("resubscribe after wait conn ready [%ds]", t)
			})
		default:
			delay.Reset()
			log.Debug("build subscribe handler")
			handleSubscribe(session, r, rxHandler)
		}
	}
}

// handleSubscribe channel断开时
func handleSubscribe(session *sessionRabbit, r Route, rxHandler RxHandler) {

	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()

	if !session.isReady {
		log.Info("session is unready")
		return
	}
	ch, ok, err := session.openChannel(func(channel *amqp.Channel) (e error) {
		e = channel.Confirm(false)

		if e != nil {
			return
		}
		_, e = channel.QueueDeclare(
			r.Queue,
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)

		if e != nil {
			return
		}
		return
	})

	if ok {
		defer func() {
			log.Debug("close channel")
			_ = ch.Close()
		}()
	}

	if err != nil {
		log.Info("initConnection channel fail")
		return
	}
	var onPublished = make(chan amqp.Confirmation)
	ch.NotifyPublish(onPublished)

	err = ch.QueueBind(r.Queue, r.Key, r.Exchange, false, nil)
	if err != nil {
		log.Info("bind queue fail: %s", r.Queue)
		log.Info(err)
		return
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Info("set qos fail")
		log.Info(err)
		return
	}

	var txHandler TxHandler = func(m Message) (err error) {

		err = publish(session, ch, m, onPublished)
		return err
	}

	msgs, err := session.stream(r.Queue, ch)
	if err != nil {
		log.Info("list message fail")
		log.Info(err)
		return
	}

	for {
		select {
		case <-session.ctx.Done(): // session销毁
			return
		case <-session.connCtx.Done(): // 连接断开
			return
		case m, ok := <-msgs:
			if !ok { // channel断开
				log.Info("message list channel closed")
				return
			}
			handleMessage(m, txHandler, rxHandler)
		}
	}
}

func handleMessage(m amqp.Delivery, txHandler TxHandler, rxHandler RxHandler) {

	var msgId = m.MessageId
	log.Info("receive message: %s", msgId)

	// TODO cache message id
	var err = rxHandler(m, txHandler)
	if err != nil {
		// TODO 从去重cache里删除message id
	} else {
		_ = m.Ack(true)
	}
}
