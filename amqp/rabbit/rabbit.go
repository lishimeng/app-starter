package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Session interface {
	Publish(m Message) error
	Close() error
	Subscribe(r Route, rxHandler RxHandler)
}

type Message struct {
	Payload interface{}
	Router  Route
	Options []PublishOption
}

func (message *Message) SetOption(option ...PublishOption) {
	message.Options = append(message.Options, option...)
}

type sessionRabbit struct {
	name            string
	ctx             context.Context
	connCtx         context.Context
	CloseConn       context.CancelFunc
	onConnClose     chan *amqp.Error
	conn            *amqp.Connection
	isReady         bool
	globalTxChannel chan Message
}

const defaultExchange = "amq.direct"

type TxHandler func(m Message) (err error)
type RxHandler func(msg amqp.Delivery, txHandler TxHandler) (err error)

type PublishOption func(m amqp.Publishing, payload interface{}) (amqp.Publishing, error)

type Route struct {
	Exchange string
	Key      string
	Queue    string
}

const (
	// When reconnecting to the server after conn failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

var (
	ErrNotConnected   = errors.New("not connected to a server")
	ErrPublishTimeout = errors.New("publish timeout")
)

var (
	JsonEncodeOption PublishOption = func(m amqp.Publishing, payload interface{}) (p amqp.Publishing, err error) {
		bs, err := json.Marshal(payload)
		if err != nil {
			return
		}
		p = m
		p.ContentType = "application/json"
		p.Body = bs
		return
	}
	TextEncodeOption PublishOption = func(m amqp.Publishing, payload interface{}) (p amqp.Publishing, err error) {
		val, ok := payload.(string)
		if !ok {
			err = fmt.Errorf("need txt payload")
			return
		}
		p = m
		p.ContentType = "text/plain"
		p.Body = []byte(val)
		return
	}
	MessageIdOption PublishOption = func(m amqp.Publishing, payload interface{}) (p amqp.Publishing, err error) {
		p = m
		p.MessageId = "1"
		return
	}
	defaultPublishOption = JsonEncodeOption
)

var MaxTxBuffer = 1024

func New(ctx context.Context, addr string) Session {

	var connCtx, cancel = context.WithCancel(context.Background())

	session := sessionRabbit{
		ctx:     ctx,
		connCtx: connCtx,
		isReady: false,
	}
	cancel() // 默认不可用
	session.initResource()

	go session.handleReconnect(addr)
	var h Session = &session
	return h
}
