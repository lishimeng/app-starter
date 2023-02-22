package app

import (
	"context"
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/app-starter/mqtt"
	shutdown "github.com/lishimeng/go-app-shutdown"
	"github.com/lishimeng/go-log"
	"testing"
	"time"
)

func TestStartMqtt(t *testing.T) {
	topic := "sss/bbb/ttt"
	var broker = "mqtt://test-mq.thingplecloud.com:1883"

	time.AfterFunc(time.Second*60, func() {
		shutdown.Exit("bye bye")
	})

	time.AfterFunc(time.Second*3, func() {
		t.Log("publish process")
		for {
			select {
			case <-factory.GetCtx().Done():
				return
			case <-time.After(time.Second):
				session := GetMqtt()
				if session != nil {
					log.Debug("pub-->")
					_ = session.Publish(topic+"1", 1, false, []byte("message content"))
					_ = session.Publish(topic+"2", 1, false, []byte("message content"))
				}
			}
		}
	})
	t.Log("start app")
	_ = New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {

		builder.
			EnableMqtt(
				mqtt.WithBroker(broker),
				mqtt.WithAuth("veolia_test_server", "f2383236")).
			ComponentAfter(func(ctx context.Context) (err error) {
				err = GetMqtt().Subscribe(func(topic string, payload []byte) {
					log.Debug("notify:%s[%s]", topic, string(payload))
				}, 1, topic+"1")

				err = GetMqtt().Subscribe(func(topic string, payload []byte) {
					log.Debug("notify:%s[%s]", topic, string(payload))
				}, 1, topic+"2")
				return
			})
		return nil
	}, func(s string) {

	})
}
