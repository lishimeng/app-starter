package client

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	ID     string              // 客户端唯一标识
	w      http.ResponseWriter // 响应写入器
	r      *http.Request
	closed bool // 关闭信号通道
}

func New(ctx context.Context, id string, request *http.Request, writer http.ResponseWriter) *Client {
	c, cancel := context.WithCancel(ctx)
	return &Client{
		ctx:    c,
		cancel: cancel,
		ID:     id,
		w:      writer,
		r:      request,
	}
}

func (c *Client) Close() {
	if c.closed {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return // lock后重复检查一遍
	}
	c.closed = true
	c.cancel()
}

func (c *Client) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
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
		// 心跳包发送逻辑
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				log.Printf("客户端 %s 心跳发送失败: %v", c.ID, err)
				c.Close()
				return
			}

		// 客户端自身Context取消（主动关闭/消息发送失败）
		case <-c.r.Context().Done():
			log.Printf("客户端 %s 自身Context取消: %v", c.ID, c.ctx.Err())
			return
		case <-c.ctx.Done():
			// 系统关闭
			return
		}
	}
}

func (c *Client) sendHeartbeat() error {
	return c.send("ping")
}

func (c *Client) SendMessage(msg string) error {
	return c.send(msg)
}

func (c *Client) send(msg string) error {
	if c.IsClosed() {
		return context.Canceled
	}

	_, err := c.w.Write([]byte("data: " + msg + "\n\n"))
	if err != nil {
		return err
	}

	if flusher, ok := c.w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}
