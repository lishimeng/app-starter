package internal

import "github.com/lishimeng/app-starter/amqp/rabbit"

type Connector struct {
	Conn string
}

// Downstream
// 监听一个数据源 Down针对broker是rx
type Downstream interface {
	Subscribe(topic string, v interface{}, txHandler rabbit.TxHandler)
	Router() rabbit.Route
}

// Upstream
// 下发数据到broker tx
type Upstream interface {

	// Submit
	// 发送到默认exchange
	Submit(topic string, v interface{}) // 提交

	// SubmitTo
	// 发送到指定exchange
	SubmitTo(exchange string, topic string, v interface{}) // 提交
}

type Service interface {
	Run() error
	UpstreamHandler() Upstream
}
