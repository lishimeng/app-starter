package rabbit

import (
	"github.com/lishimeng/go-log"
	"github.com/streadway/amqp"
	"time"
)

func subscribe(session *sessionRabbit, r Route, rxHandler RxHandler) {

	if len(r.Queue) <= 0 {
		log.Info("queue is empty, exit")
		return
	}
	if len(r.Key) <= 0 {
		r.Key = r.Queue // 如果不提供Key,使用queue名字作为key
	}
	if len(r.Exchange) <= 0 {
		r.Exchange = defaultExchange
	}

	var delay = NewDelay(1, 60, false)

	for {
		select {
		case <-session.ctx.Done(): // session销毁
			log.Fine("subscribe exit [%s:%s-%s]", r.Exchange, r.Key, r.Queue)
			return
		case <-session.connCtx.Done(): // 连接断开
			// 等待连接
			delay.Delay(func(t int) {
				log.Fine("resubscribe after wait conn ready [%ds]", t)
			})
		default:
			delay.Reset()
			log.Fine("build subscribe handler")
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

	var serverCtx = ServerContext{Router: r}

	if !session.isReady {
		log.Debug("session is unready")
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

		err = session.Publish(m)
		return err
	}

	msgQueue, err := session.stream(r.Queue, ch)
	if err != nil {
		log.Info("list message fail")
		log.Info(err)
		return
	}

	go func() {
		for {
			select {
			case <-session.ctx.Done():
				return
			case <-session.connCtx.Done():
				return
			case <-time.After(time.Second * 10):
				q, e := ch.QueueInspect(r.Queue)
				if e != nil {
					return
				}
				serverCtx.Messages = q.Messages
				serverCtx.Consumers = q.Consumers
			}
		}
	}()

	for {
		select {
		case <-session.ctx.Done(): // session销毁
			return
		case <-session.connCtx.Done(): // 连接断开
			return
		case m, ok := <-msgQueue:
			if !ok { // channel断开
				log.Info("message list channel closed")
				return
			}
			handleMessage(m, txHandler, rxHandler, serverCtx)
		}
	}
}

func handleMessage(m amqp.Delivery, txHandler TxHandler, rxHandler RxHandler, ctx ServerContext) {

	var msgId = m.MessageId
	log.Fine("receive message: %s", msgId)

	// TODO cache message id
	var err = rxHandler(m, txHandler, ctx)
	if err != nil {
		// TODO 从去重cache里删除message id
		e := m.Nack(false, true)
		if e != nil {
			log.Info("nack fail")
			log.Info(e)
		}
	} else {
		e := m.Ack(true)
		if e != nil {
			log.Info("ack fail")
			log.Info(e)
		}
	}
}
