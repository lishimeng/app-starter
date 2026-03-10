package client

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client 表示一个SSE客户端连接
type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
	ID     string              // 客户端唯一标识
	w      http.ResponseWriter // 响应写入器
	r      *http.Request
	closed bool            // 关闭信号
	events map[string]bool // 订阅的事件集合
}

func New(ctx context.Context, id string, request *http.Request, writer http.ResponseWriter) *Client {
	c, cancel := context.WithCancel(ctx)
	return &Client{
		ctx:    c,
		cancel: cancel,
		ID:     id,
		w:      writer,
		r:      request,
		events: make(map[string]bool),
	}
}

func (c *Client) Close() {
	if c.closed {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.cancel()
}

func (c *Client) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Subscribe 订阅指定事件
func (c *Client) Subscribe(events ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range events {
		c.events[e] = true
	}
}

// Unsubscribe 取消订阅指定事件
func (c *Client) Unsubscribe(events ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range events {
		delete(c.events, e)
	}
}

// IsSubscribed 检查是否订阅了指定事件
func (c *Client) IsSubscribed(event string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.events[event]
}

// GetEvents 获取所有订阅的事件列表
func (c *Client) GetEvents() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]string, 0, len(c.events))
	for e := range c.events {
		result = append(result, e)
	}
	return result
}

func (c *Client) Run(heartbeatInterval time.Duration) {
	if c.IsClosed() {
		log.Printf("客户端 %s 已关闭，跳过主循环", c.ID)
		return
	}

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	log.Printf("客户端 %s 启动主循环", c.ID)
	for {
		select {
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				log.Printf("客户端 %s 心跳发送失败: %v", c.ID, err)
				c.Close()
				return
			}

		case <-c.r.Context().Done():
			log.Printf("客户端 %s 自身Context取消: %v", c.ID, c.ctx.Err())
			return
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) sendHeartbeat() error {
	return c.send("data: ping\n\n")
}

func (c *Client) SendMessage(msg string) error {
	return c.send(msg)
}

func (c *Client) send(msg string) error {
	if c.IsClosed() {
		return context.Canceled
	}

	_, err := c.w.Write([]byte(msg))
	if err != nil {
		return err
	}

	if flusher, ok := c.w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
