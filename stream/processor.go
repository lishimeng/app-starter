package stream

import "io"

func (s *SessionCtx) rxLoop() {
	var buf = make([]byte, 1024)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			n, err := s.proxy.Read(buf)
			if err != nil {
				if err == io.EOF {
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
		_ = s.packetBuf.Next(n)
	}
}
