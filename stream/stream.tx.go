package stream

import (
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
)

func (s *SessionCtx) Write(p []byte) (size int, err error) {
	if ioLog {
		log.Info(">>%s", tool.BytesToHex(p))
	}

	size, err = s.proxy.Write(p)
	return
}
