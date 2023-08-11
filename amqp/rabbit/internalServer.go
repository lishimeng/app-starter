package rabbit

import (
	"errors"
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func (session *sessionRabbit) handleReconnect(addr string) {
	for {
		session.isReady = false
		log.Fine("Attempting to connect")

		conn, err := session.connect(addr)

		if err != nil {
			log.Debug("Failed to connect. Retrying...")

			select {
			case <-session.ctx.Done():
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		// after connect to broker success, initConnection the session
		if done := session.sessionLoop(conn); done {
			break
		} else {
			log.Debug("reconnect")
		}
	}
}

// sessionLoop will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *sessionRabbit) sessionLoop(conn *amqp.Connection) bool {
	for {
		session.isReady = false

		err := session.initConnection(conn)

		if err != nil {
			log.Debug("Failed to initialize conn. Retrying...")

			select {
			case <-session.ctx.Done():
				go func() {
					err = session.Close()
					if err != nil {
						log.Info(err)
					}
				}()
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.ctx.Done():
			go func() {
				err = session.Close()
				if err != nil {
					log.Info(err)
				}
			}()
			return true
		case <-session.onConnClose:
			session.isReady = false
			log.Debug("on conn closed")
			session.CloseConn()
			session.safeCloseConnection()
			return false
		}
	}
}

func (session *sessionRabbit) openChannel(afterChannel ...func(*amqp.Channel) error) (ch *amqp.Channel, ok bool, err error) {
	ch, err = session.conn.Channel()

	if err != nil {
		return
	}
	ok = true
	for _, after := range afterChannel {
		err = after(ch)
		if err != nil {
			break
		}
	}

	return
}

// Push will push data onto the queue, and wait for a confirm_signal.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirm_signal is received.
// This will block until the server sends a confirm_signal. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (session *sessionRabbit) Push(data []byte, channel *amqp.Channel) error {
	if !session.isReady {
		return errors.New("failed to push push: not connected")
	}
	for {
		err := session.UnsafePush(data, channel)
		if err != nil {
			log.Info("Push failed. Retrying...")
			select {
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		//case confirm := <-session.notifyConfirm:
		//	if confirm.Ack {
		//		log.Info("Push confirmed!")
		//		return nil
		//	}
		case <-time.After(resendDelay):
		}
		log.Info("Push didn't confirm. Retrying...")
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (session *sessionRabbit) UnsafePush(data []byte, channel *amqp.Channel) error {
	if !session.isReady {
		return ErrNotConnected
	}
	return channel.Publish(
		"",    // Exchange
		"",    // Routing key
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (session *sessionRabbit) stream(subscribe string, channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	if !session.isReady {
		return nil, ErrNotConnected
	}
	return channel.Consume(subscribe, "", false, false, false, false, nil)
}

// disposeResources 清理资源
func (session *sessionRabbit) disposeResources() error {

	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()

	var err error = nil

	session.releaseResource()
	return err
}

func (session *sessionRabbit) monitor() {
	go func() {
		for {
			select {
			case <-session.ctx.Done():
				return
			case <-time.After(time.Second * 5):
				if len(session.globalTxChannel) <= 0 {
					break
				}
				log.Fine("tx status:%d[%d]", len(session.globalTxChannel), cap(session.globalTxChannel))
			}
		}
	}()
}
