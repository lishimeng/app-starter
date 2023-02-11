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

func (s *simpleDs) Subscribe(_ interface{}, _ rabbit.TxHandler, serverContext rabbit.ServerContext) {
	rxIndex++
	log.Info("receive:%d", rxIndex)
	log.Info("%+v", serverContext)
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
	rabbit.MaxTxBuffer = 2048

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
						if i%5 == 0 {
							m.SetOption(rabbit.TextEncodeOption, rabbit.UUIDMsgIdOption)
						} else {
							m.SetOption(rabbit.JsonEncodeOption, rabbit.UUIDMsgIdOption)
						}

						e := Publish(session, m)
						if e != nil {
							log.Info("publish timeout")
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
