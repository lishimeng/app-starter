package stream

import (
	"context"
	"github.com/tarm/serial"
	"time"
)

func NewSerialSession(ctx context.Context, name string, baud int, opts ...Option) (s Session, err error) {
	c := &serial.Config{
		Name:        name,
		Baud:        baud,
		Size:        8,
		ReadTimeout: time.Second * 5,
		StopBits:    1,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done(): // 正常关闭
				if port == nil {
					return
				}
				_ = port.Close()
			}
		}
	}()
	s = NewSession(ctx, port, opts...)
	return
}
