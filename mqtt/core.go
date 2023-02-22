package mqtt

import (
	"context"
	"errors"
	proxy "github.com/eclipse/paho.mqtt.golang"
	"github.com/lishimeng/go-log"
	"sync"
	"time"
)

type sessionContext struct {
	conn   proxy.Client
	locker *sync.Mutex
	ctx    context.Context
}

func New(ctx context.Context, options ...ClientOption) Session {
	opt := proxy.NewClientOptions()
	for _, option := range options {
		opt = option(opt)
	}

	client := proxy.NewClient(opt)
	var session = sessionContext{
		conn:   client,
		ctx:    ctx,
		locker: &sync.Mutex{},
	}
	go session.listenExit()
	return &session
}

func (session *sessionContext) listenExit() {
	for {
		select {
		case <-session.ctx.Done():
			log.Debug("mqtt exit")
			if session.conn != nil && session.conn.IsConnectionOpen() {
				session.conn.Disconnect(10)
			}
			return
		}
	}
}

func (session *sessionContext) Connect() error {
	return session.ensureConnected()
}

func (session *sessionContext) ensureConnected() error {
	if !session.conn.IsConnected() {
		session.locker.Lock()
		defer session.locker.Unlock()
		if !session.conn.IsConnected() {
			if token := session.conn.Connect(); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}
	return nil
}

func (session *sessionContext) Publish(topic string, qos byte, retained bool, data []byte) error {
	if err := session.ensureConnected(); err != nil {
		return err
	}

	token := session.conn.Publish(topic, qos, retained, data)
	if err := token.Error(); err != nil {
		return err
	}

	// return false is the timeout occurred
	if !token.WaitTimeout(time.Second * 10) {
		return ErrPublishTimeout
	}

	return nil
}

func (session *sessionContext) Subscribe(handler func(topic string, payload []byte), qos byte, topic string) error {
	if len(topic) == 0 {
		return errors.New("the topic is empty")
	}

	session.conn.Subscribe(topic, qos, func(_ proxy.Client, message proxy.Message) {
		if message.Duplicate() {
			return
		}
		handleMessage(message, handler)
	})
	return nil
}

func handleMessage(message proxy.Message, handler func(topic string, payload []byte)) {
	handler(message.Topic(), message.Payload())
	return
}

func (session *sessionContext) Unsubscribe(topics ...string) {
	session.conn.Unsubscribe(topics...)
}
