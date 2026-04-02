package event

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType int

const (
	ToClient  MessageType = 1
	Broadcast MessageType = 2
)

type Payload struct {
	Data  json.RawMessage `json:"data"`  // 业务数据（JSON格式）
	Event string          `json:"event"` // SSE事件类型（如：order_update、new_notification）
	ID    string          `json:"id"`    // 事件ID（可选）
	Time  int64           `json:"time"`  // 事件时间戳(deprecate)
	Retry int64           `json:"retry"` // ms
}

type Event struct {
	Type     MessageType // 发送类型
	ClientId string      // 发送到指定客户端时需要提供
	Payload  *Payload
}

func (e *Payload) Marshall() string {
	// 填充默认值
	if e.Time == 0 {
		e.Time = time.Now().UnixMilli()
	}
	if e.Event == "" {
		e.Event = "message"
	}

	// 按SSE协议拼接（支持自定义event类型）
	// 格式：event: {event}\nid: {id}\ndata: {json}\n\n
	sseStr := ""
	if len(e.Event) > 0 { // 可选
		sseStr += fmt.Sprintf("event: %s\n", e.Event)
	}
	if len(e.ID) > 0 { // 可选
		sseStr += fmt.Sprintf("id: %s\n", e.ID)
	}
	if e.Retry > 0 { // 可选
		sseStr += fmt.Sprintf("retry: %d\n", e.Retry)
	}
	if e.Time > 0 { // 可选
		sseStr += fmt.Sprintf("time: %d\n", e.Time)
	}
	sseStr += "data: " + string(e.Data) + "\n\n" // 必选
	return sseStr
}

func New(eventType string, data any) (*Payload, error) {
	// 序列化业务数据为JSON
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Payload{
		Data:  dataJson,
		Event: eventType,
		Time:  time.Now().UnixMilli(),
	}, nil
}
