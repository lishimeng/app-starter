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

func (s *simpleDs) Subscribe(_ string, _ interface{}, _ rabbit.TxHandler) {
	rxIndex++
	log.Info("receive:%d", rxIndex)
}

func (s *simpleDs) Router() rabbit.Route {
	return rabbit.Route{
		Exchange: "",
		Key:      "wwwwwww",
		Queue:    "qqqqqqqqq",
	}
}

func TestSdk001(t *testing.T) {

	log.SetLevelAll(log.DEBUG)

	const addr = "amqp://ows:thingple@127.0.0.1:5672/"

	var ctx, cancel = context.WithCancel(context.Background())
	var c = Connector{Conn: addr}
	var session = New(ctx, c)

	log.Info(session)

	var ds Handler = &simpleDs{}

	log.Info("register subscriber")
	go RegisterHandler(session, ds)

	time.Sleep(time.Second * 3)
	go func() {

		for {
			r := rand.Intn(10) + 5
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(r)):
				go func() {
					log.Info("publish:---> %d", r)
					for i := 0; i < r; i++ {
						var txt = fmt.Sprintf("message %d", i+1)
						var m = rabbit.Message{
							Payload: []byte(txt),
							Router:  ds.Router(),
						}
						m.SetOption(rabbit.JsonEncodeOption, rabbit.MessageIdOption)
						Publish(session, m)
					}
				}()
			}
		}

	}()

	time.Sleep(time.Minute * 20)
	log.Info("done")
	cancel()

	time.Sleep(time.Second * 3)
	log.Info("close")
}
