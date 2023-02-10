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

func RegisterHandler(session rabbit.Session, ds ...Downstream) {
	for _, d := range ds {
		registerHandler(session, d)
	}
}

func registerHandler(session rabbit.Session, ds Downstream) {
	session.Subscribe(ds.Router(), func(msg amqp.Delivery, txHandler rabbit.TxHandler) (err error) {
		ds.Subscribe(ds.Router().Queue, msg.Body, txHandler)
		return nil
	})
}

func Publish(session rabbit.Session, m rabbit.Message) {
	session.Publish(m)
}
