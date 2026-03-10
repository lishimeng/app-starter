package event

import (
	"encoding/json"
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
	Time  int64           `json:"time"`  // 事件时间戳
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

	// 序列化为JSON字符串
	jsonData, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	// 按SSE协议拼接（支持自定义event类型）
	// 格式：event: {event}\nid: {id}\ndata: {json}\n\n
	sseStr := ""
	if e.ID != "" {
		sseStr += "id: " + e.ID + "\n"
	}
	sseStr += "event: " + e.Event + "\n"
	sseStr += "data: " + string(jsonData) + "\n\n"
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
