package amqp

import (
	"fmt"
	"github.com/lishimeng/app-starter/amqp/rabbit"
)

const addrFormat = "amqp://%s:%s@%s:%d/"

type Connector struct {
	Conn string
}

func (c *Connector) Build(host string, port int, user string, passwd string) {
	c.Conn = fmt.Sprintf(addrFormat, user, passwd, host, port)
}

// Handler
// 监听一个数据源 Down针对broker是rx
type Handler interface {
	Subscribe(topic string, v interface{}, txHandler rabbit.TxHandler)
	Router() rabbit.Route
}
