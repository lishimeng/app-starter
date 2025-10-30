package amqp

import (
	"context"
	"errors"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connTpl = "amqp://%s:%s@%s:%d/"

var (
	ErrEmptyMessageRouter = errors.New("empty router key")
	ErrSessionNil         = errors.New("session is nil")
)

func New(ctx context.Context, c Connector, options ...rabbit.SessionOption) (session rabbit.Session) {
	session = rabbit.New(ctx, c.conn, options...)
	return
}

func RegisterHandler(session rabbit.Session, handlers ...Handler) {
	for _, handler := range handlers {
		registerHandler(session, handler)
	}
}

func registerHandler(session rabbit.Session, handler Handler) {
	if session == nil || handler == nil {
		return
	}

	r := handler.Router()
	if len(r.Queue) <= 0 {
		log.Info("queue is empty, exit")
		return
	}
	if len(r.Key) <= 0 {
		r.Key = r.Queue // 如果不提供Key,使用queue名字作为key
	}
	if len(r.Exchange) <= 0 {
		r.Exchange = rabbit.DefaultExchange
	}
	session.Subscribe(r,
		func(msg amqp.Delivery, txHandler rabbit.TxHandler, serverCxt rabbit.ServerContext) (err error) {
			handler.Subscribe(msg.Body, txHandler, serverCxt)
			return
		})
}

// Publish 发送buffer满之后,返回rabbit.ErrPublishTimeout
func Publish(session rabbit.Session, m rabbit.Message) error {
	if session == nil {
		return ErrSessionNil
	}
	if m.Payload == nil {
		// payload empty
		return nil
	}
	if len(m.Router.Exchange) <= 0 {
		log.Debug("use default exchange:%s", rabbit.DefaultExchange)
		m.Router.Exchange = rabbit.DefaultExchange
	}
	if len(m.Router.Key) <= 0 {
		log.Debug("use queue as router key:%s", m.Router.Queue)
		m.Router.Key = m.Router.Queue
	}

	if len(m.Router.Key) <= 0 {
		// 无目的地, 放弃publish操作
		log.Debug("empty router key, drop this message")
		return ErrEmptyMessageRouter
	}
	return session.Publish(m)
}
