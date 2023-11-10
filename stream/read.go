package stream

import (
	"github.com/lishimeng/go-log"
	"io"
)

func (s *SessionCtx) OnData(p []byte, size int, err error) {
	if s.readSubscriber == nil {
		return
	}
	var payload = make([]byte, size)
	copy(payload, p[:size])
	s.readSubscriber(payload, err)
}

func (s *SessionCtx) Read(p []byte) (size int, err error) {
	size, err = s.proxy.Read(p)
	if err != nil {
		if err != io.EOF { // 关闭了
			log.Info(err)
			panic(err)
		}
		err = nil // 吃掉err
	}
	return
}
