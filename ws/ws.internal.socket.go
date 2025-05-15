package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lishimeng/app-starter/broker"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/go-log"
	"strings"
	"time"
)

type Downstream struct {
	Channel  string `json:"channel,omitempty"`  // 接收者通道
	Category string `json:"category,omitempty"` // 类型: mq/api
	Payload  any    `json:"payload,omitempty"`  // 数据
}

//type TxHandler func(m any)

type TxHandler func(channel string, category string, payload any)

func handler(ctx context.Context, topic string, ws *websocket.Conn, logic Logic) {
	var queue = make(chan any)
	var txQueue = make(chan *websocket.PreparedMessage, 100)
	var txBuf bytes.Buffer
	var txEncoder = json.NewEncoder(&txBuf)

	var apiReq = make(chan RestJob, 32)

	var jsonTxFunc = func(channel string, category string, payload any) {
		var obj = Downstream{
			Channel:  channel,
			Category: category,
			Payload:  payload,
		}
		if e := txEncoder.Encode(obj); e != nil {
			log.Info(e)
		} else {
			content := txBuf.String()
			txBuf.Reset()
			msg, e := websocket.NewPreparedMessage(websocket.TextMessage, []byte(content))
			if e != nil {
				log.Info(e)
				return
			}
			txQueue <- msg
		}
	}

	go apiProxyWatchDog(ctx, apiReq, jsonTxFunc)

	exit := broker.Get().Subscribe(topic, func(msg broker.MessageItem) {
		m := msg.Message
		queue <- m
	})
	defer func() {
		exit()
	}()

	// tx loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-txQueue:
				if !ok {
					return
				}
				_writeMsg(ws, m)
			}
		}
	}()

	// ping loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 3):
				msg, e := websocket.NewPreparedMessage(websocket.PingMessage, nil)
				if e != nil {
					log.Info(e)
					return
				}
				txQueue <- msg
			}
		}
	}()

	var messageHandler = func(m []byte) {
		// 处理api请求
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		// 转换成api req
		var req RestJob
		err := json.Unmarshal(m, &req)
		if err != nil {
			log.Info(err)
			log.Info(string(m))
			return
		}
		apiReq <- req
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				messageType, payload, err := ws.ReadMessage()
				if err != nil {
					log.Info(err)
					exit()
					return
				}
				switch messageType {
				case websocket.TextMessage:
					go messageHandler(payload)
				}
			}
		}

	}()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-queue:
			if !ok {
				return
			}
			logic(data, jsonTxFunc)
		}
	}
}

func handleWsSession(ctx server.Context, topic string, logic Logic) {
	sessionId := strings.ReplaceAll(uuid.New().String(), "-", "")
	log.Info("create WS Session: %s", sessionId)
	ws, err := upgrade.Upgrade(ctx.C.ResponseWriter(), ctx.C.Request(), nil)
	if err != nil {
		var handshakeError websocket.HandshakeError
		if !errors.As(err, &handshakeError) {
			log.Info(err)
		}
		return
	}
	wait, cancel := context.WithCancel(context.Background())
	ws.SetCloseHandler(func(code int, text string) error {
		log.Info("WS Session close: %s", sessionId)
		cancel()
		return nil
	})

	defer func() {
		_ = ws.Close()
	}()

	handler(wait, topic, ws, logic)
}

func _writeMsg(conn *websocket.Conn, m *websocket.PreparedMessage) {
	defer func() {
		if e := recover(); e != nil {
			log.Info(e)
		}
	}()
	_ = conn.WritePreparedMessage(m)
}
