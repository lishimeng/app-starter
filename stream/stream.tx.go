package stream

import (
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/util"
)

func (s *SessionCtx) Write(p []byte) (size int, err error) {
	if ioLog {
		log.Info(">>%s", util.BytesToHex(p))
	}

	size, err = s.proxy.Write(p)
	return
}
