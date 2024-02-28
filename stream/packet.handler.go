package stream

import (
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/util"
	"io"
)

func (s *SessionCtx) rxLoop() {
	defer func() {
		if e := recover(); e != nil {
			if s.onErr != nil {
				s.onErr(e)
			}
		}
	}()
	var buf = make([]byte, 1024)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			n, err := s.proxy.Read(buf)
			if err != nil {
				if err == io.EOF { // 关闭了
					panic(err)
				}
			}
			s.packetBuf.Write(buf[:n])
			s._packet()
		}
	}
}

func (s *SessionCtx) _packet() {
	if s.pp == nil {
		return
	}
	n := s.pp.Data([]byte(s.packetBuf.String()))
	if n > 0 {
		packet := s.packetBuf.Next(n)
		if ioLog {
			log.Info("<<%s", util.BytesToHex(packet))
		}
	}
}
