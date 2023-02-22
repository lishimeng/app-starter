package mqtt

import (
	"context"
	"github.com/lishimeng/go-log"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	var ctx, exit = context.WithCancel(context.Background())
	var err error
	var broker = "mqtt://127.0.0.1:1883"
	var topic = "p_test_mqtt_lib"
	var qos byte = 0
	var session = New(ctx, WithBroker(broker),
		WithRandomClientId(),
		WithAuth("veolia_test_server", "f2383236"))

	err = session.Connect()
	if err != nil {
		t.Fatal(err)
	}

	time.AfterFunc(time.Second*30, func() {
		log.Info("exit")
		exit()
	})

	err = session.Subscribe(func(topic string, payload []byte) {
		log.Debug("receive1:%s[%s]", topic, string(payload))
	}, qos, topic)

	err = session.Subscribe(func(topic string, payload []byte) {
		log.Debug("receive2:%s[%s]", topic, string(payload))
	}, qos, topic+"2")

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				log.Info("publish a message")
				err = session.Publish(topic, qos, false, []byte("test mqtt lib"))
				if err != nil {
					log.Debug(err)
				}
				err = session.Publish(topic+"2", qos, false, []byte("test mqtt lib222"))
				if err != nil {
					log.Debug(err)
				}
			}
		}
	}()

	<-ctx.Done()
}
