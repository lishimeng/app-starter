package rabbit

import (
	"github.com/lishimeng/go-log"
)

func (session *sessionRabbit) initResource() {
	session.globalTxChannel = make(chan Message, MaxTxBuffer)
}

func (session *sessionRabbit) releaseResource() {
	session.safeCloseConnection()
}

func (session *sessionRabbit) safeCloseConnection() {
	log.Fine("safe close conn")
	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()
	var err error = nil
	if session.conn != nil {
		var conn = session.conn
		session.conn = nil
		err = conn.Close()
	}
	if err != nil {
		log.Debug(err)
	}

}
