package rabbit

import (
	"github.com/lishimeng/go-log"
)

func (session *sessionRabbit) initResource() {
	//session.onConnClose = make(chan *amqp.Error)
	session.globalTxChannel = make(chan Message, MaxTxBuffer)
}

func (session *sessionRabbit) releaseResource() {
	session.safeCloseGlobalTx()
	session.safeCloseConnection()
	//session.safeCloseConnCloseChan()
}

func (session *sessionRabbit) safeCloseConnection() {
	log.Debug("safe close conn")
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

func (session *sessionRabbit) safeCloseGlobalTx() {
	log.Debug("safe close tx")
	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()
	close(session.globalTxChannel)
}
func (session *sessionRabbit) safeCloseConnCloseChan() {
	log.Debug("safe close conn notify")
	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()
	close(session.onConnClose)
}
