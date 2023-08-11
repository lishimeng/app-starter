package rabbit

import (
	"context"
	"fmt"
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// connect will create a new AMQP conn
func (session *sessionRabbit) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	session.changeConnection(conn)
	log.Info("Amqp Connected!")
	return conn, nil
}

// initConnection will initialize channel & declare queue
func (session *sessionRabbit) initConnection(_ *amqp.Connection) error {
	session.isReady = true
	for i := 1; i <= session.multiTx; i++ {
		name := fmt.Sprintf("tx_process_%d", i)
		go session.txProcess(name)
	}
	return nil
}

// changeConnection takes a new conn to the queue,
// and updates the close listener to reflect this.
func (session *sessionRabbit) changeConnection(connection *amqp.Connection) {
	session.conn = connection
	session.onConnClose = make(chan *amqp.Error)
	session.conn.NotifyClose(session.onConnClose)
	session.connCtx, session.CloseConn = context.WithCancel(session.ctx)
}
