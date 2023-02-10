package amqp

import (
	"context"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/streadway/amqp"
)

func New(ctx context.Context, c Connector) (session rabbit.Session) {
	session = rabbit.New(ctx, c.Conn)
	return
}

func RegisterHandler(session rabbit.Session, handlers ...Handler) {
	for _, handler := range handlers {
		registerHandler(session, handler)
	}
}

func registerHandler(session rabbit.Session, handler Handler) {
	session.Subscribe(handler.Router(), func(msg amqp.Delivery, txHandler rabbit.TxHandler) (err error) {
		handler.Subscribe(handler.Router().Queue, msg.Body, txHandler)
		return nil
	})
}

// Publish 发送buffer满之后,返回rabbit.ErrPublishTimeout
func Publish(session rabbit.Session, m rabbit.Message) error {
	return session.Publish(m)
}
