package ws

import (
	"context"
	"github.com/lishimeng/go-log"
	"time"
)

func apiProxyWatchDog(ctx context.Context, req chan RestJob, respCallback TxHandler) {
	c, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			apiProxyLoop(c, req, respCallback)
		}
		time.Sleep(time.Millisecond * 1000) // 如果是异常, 暂停一段时间
	}
}

func apiProxyLoop(ctx context.Context, req chan RestJob, respCallback TxHandler) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
			return
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-req:
			if !ok {
				return
			}
			apiProxyHandler(d, respCallback)
		}
	}

}

func apiProxyHandler(req RestJob, respCallback TxHandler) {
	var resp = make(map[string]any)
	err := req.Fetch(&resp)
	if err != nil {
		log.Info(err)
	}
	respCallback(req.Api, "api", resp) // TODO
	return
}
