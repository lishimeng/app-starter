package amqp

import (
	"context"
	"fmt"
	"github.com/lishimeng/app-starter/amqp/rabbit"
	"github.com/lishimeng/go-log"
	"math/rand"
	"testing"
	"time"
)

type simpleDs struct {
}

var rxIndex = 0

const (
	testUser   = "office"
	testPasswd = "thingple"
	testHost   = "192.168.10.254"
	testPort   = 15672
)

func (s *simpleDs) Subscribe(_ interface{}, _ rabbit.TxHandler, serverContext rabbit.ServerContext) {
	rxIndex++
	log.Info("receive:%d", rxIndex)
	log.Info("%+v", serverContext)
}

func (s *simpleDs) Router() rabbit.Route {
	return rabbit.Route{
		Exchange: "",
		Key:      "awesome_demo",
		Queue:    "awesome_demo_queue",
	}
}

func TestSdk001(t *testing.T) {

	log.SetLevelAll(log.FINE)

	var addr = fmt.Sprintf(connTpl, testUser, testPasswd, testHost, testPort)
	log.Info(addr)
	rabbit.MaxTxBuffer = 20

	var ctx, cancel = context.WithCancel(context.Background())
	var c = Connector{conn: addr}
	var session = New(ctx, c, rabbit.TxWorkerOption(3))

	log.Info(session)

	var ds Handler = &simpleDs{}

	log.Info("register subscriber")
	go RegisterHandler(session, ds)

	time.Sleep(time.Second * 3)
	go func() {

		var index = 0
		for {
			r := rand.Intn(10) + 5
			delay := 1

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(delay)):
				go func() {
					log.Info("publish:---> %d", r)

					for i := 0; i < r; i++ {
						index++
						var payload = make(map[string]interface{})
						payload["index"] = index
						payload["title"] = "sdk test"
						var m = rabbit.Message{
							Payload: payload,
							Router:  ds.Router(),
						}
						if i%5 == 0 {
							m.SetOption(rabbit.TextEncodeOption, rabbit.UUIDMsgIdOption)
						} else {
							m.SetOption(rabbit.JsonEncodeOption, rabbit.UUIDMsgIdOption)
						}

						e := Publish(session, m)
						if e != nil {
							log.Info("publish timeout")
							log.Info(e)
							log.Info("cancel")
							cancel()
						}
					}
				}()
			}
		}

	}()

	time.Sleep(time.Minute * 2)
	log.Info("done")
	cancel()

	time.Sleep(time.Second * 3)
	log.Info("close")
}

func TestPubOnce(t *testing.T) {
	log.SetLevelAll(log.FINE)

	var addr = fmt.Sprintf(connTpl, testUser, testPasswd, testHost, testPort)
	rabbit.MaxTxBuffer = 20

	var ctx, cancel = context.WithCancel(context.Background())
	var c = Connector{conn: addr}
	var session = New(ctx, c)

	log.Info(session)

	var ds Handler = &simpleDs{}

	log.Info("register subscriber")
	go func() {

		var index = 0
		for {
			r := rand.Intn(10) + 5
			delay := 1

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(delay)):
				go func() {
					log.Info("publish:---> %d", r)

					for i := 0; i < 2; i++ {
						index++
						var payload = make(map[string]interface{})
						payload["index"] = index
						payload["title"] = "sdk test"
						var m = rabbit.Message{
							Payload: payload,
							Router:  ds.Router(),
						}
						m.SetOption(rabbit.JsonEncodeOption, rabbit.UUIDMsgIdOption)

						e := Publish(session, m)
						if e != nil {
							log.Info("publish timeout")
							log.Info(e)
						}
					}
				}()
				return
			}
		}

	}()

	time.Sleep(time.Second * 60)
	log.Info("done")
	cancel()

	time.Sleep(time.Second * 3)
	log.Info("close")
}
