package broker

import (
	"context"
	"fmt"
	"github.com/lishimeng/go-log"
	"math/rand/v2"
	"sync"
	"time"
)

type OnDataFunc func(MessageItem)

// Unsubscribe 取消订阅, 支持对同一topic的多个订阅
type Unsubscribe func()

// Router 支持对同一topic的多个订阅
type Router struct {
	lock        sync.Mutex
	subscribers map[string]map[string]OnDataFunc // topic/id(callback)
}

type MessageItem struct {
	Message any
	Topic   string
}

func (r *Router) Add(id string, topic string, f OnDataFunc) Unsubscribe {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.subscribers == nil {
		r.subscribers = make(map[string]map[string]OnDataFunc)
	}
	group, ok := r.subscribers[topic]
	if !ok {
		group = make(map[string]OnDataFunc)
		r.subscribers[topic] = group
	}
	group[id] = f
	return func() {
		r.Remove(topic, id)
	}
}

func (r *Router) Remove(topic, id string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.subscribers[topic], id)
}

func (r *Router) Group(topic string) map[string]OnDataFunc {
	return r.subscribers[topic]
}

type Client struct {
	ctx     context.Context
	channel chan MessageItem
	router  *Router
}

func New(ctx context.Context) *Client {
	r := Router{}
	channel := make(chan MessageItem, 10) // 放大
	client := &Client{ctx: ctx, router: &r, channel: channel}
	go client.transmission()
	return client
}
func (c *Client) transmission() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c._transmissionLoop()
		}
	}
}

func (c *Client) _transmissionLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.channel:
			if !ok {
				return
			}
			c.broadcast(msg)
		}
	}
}

func (c *Client) broadcast(item MessageItem) {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
	}()

	group := c.router.Group(item.Topic)
	if group == nil { // 无订阅
		return
	}
	for _, callback := range group {
		go c.broadcastTo(callback, item)
	}
}

func (c *Client) broadcastTo(sendTo OnDataFunc, item MessageItem) {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
	}()
	sendTo(item)
}

func (c *Client) Publish(topic string, data any) {

	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
	}()
	ctx, cancel := context.WithTimeout(c.ctx, time.Second)
	defer func() {
		cancel()
	}()
	select {
	case <-ctx.Done():
		return
	default:
		c.channel <- MessageItem{Message: data, Topic: topic}
	}
}

func (c *Client) Subscribe(topic string, h OnDataFunc) Unsubscribe {
	if h == nil {
		return func() {}
	}
	id := fmt.Sprintf("%d.%02d", time.Now().Nanosecond(), rand.IntN(32))
	unsubscribe := c.router.Add(id, topic, h)
	return unsubscribe
}
