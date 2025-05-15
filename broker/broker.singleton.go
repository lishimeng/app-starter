package broker

import "context"

// 用单例模式创建和使用

var ins *Client

func Init(ctx context.Context) {
	ins = New(ctx)
}

func Get() *Client {
	return ins
}
