package sse

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/lishimeng/app-starter/sse/event"
	"github.com/lishimeng/go-log"
)

//go:embed static
var staticFiles embed.FS

func TestSse(t *testing.T) {

	var ctx = context.Background()
	var m = NewManager(ctx)
	m.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		eventsParam := r.URL.Query().Get("events")
		var events []string
		if eventsParam != "" {
			events = strings.Split(eventsParam, ",")
			for i := range events {
				events[i] = strings.TrimSpace(events[i])
			}
		}
		WebHandler(m, w, r, events...)
	})
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Info("提取static目录失败:", err)
	}

	pages := http.FS(staticFS)

	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		log.Info("view page: /")
		f, e := pages.Open("page_view.html")
		if e != nil {
			log.Info(err)
			w.WriteHeader(500)
			return
		}
		defer f.Close()
		data, e := io.ReadAll(f)
		if e != nil {
			log.Info(err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	go func() {
		var i = 0
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				i++
				log.Info("server[broadcast]--> all client %d", i)
				payload, _ := event.New("system_msg", map[string]any{
					"msg": fmt.Sprintf("broadcast msg: %d", i),
				})
				evt := event.Event{
					Type:    event.Broadcast,
					Payload: payload,
				}
				_ = m.SendEvent(evt)
			}
		}
	}()

	log.Info("startup web server...")
	_ = http.ListenAndServe(":8080", mux)
}

func TestSseWithChannel(t *testing.T) {
	var ctx = context.Background()
	var m = NewManager(ctx)
	m.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		eventsParam := r.URL.Query().Get("events")
		var events []string
		if eventsParam != "" {
			events = strings.Split(eventsParam, ",")
			for i := range events {
				events[i] = strings.TrimSpace(events[i])
			}
		}
		WebHandler(m, w, r, events...)
	})

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Info("提取static目录失败:", err)
	}

	pages := http.FS(staticFS)

	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		log.Info("view page: /")
		f, e := pages.Open("page_view.html")
		if e != nil {
			log.Info(err)
			w.WriteHeader(500)
			return
		}
		defer f.Close()
		data, e := io.ReadAll(f)
		if e != nil {
			log.Info(err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	go func() {
		orderCount := 0
		noticeCount := 0
		systemCount := 0

		tickerOrder := time.NewTicker(time.Second * 5)
		tickerNotice := time.NewTicker(time.Second * 7)
		tickerSystem := time.NewTicker(time.Second * 12)
		defer tickerOrder.Stop()
		defer tickerNotice.Stop()
		defer tickerSystem.Stop()

		for {
			select {
			case <-tickerOrder.C:
				orderCount++
				payload, _ := event.New("order", map[string]any{
					"orderId": fmt.Sprintf("ORD-%d", orderCount),
					"status":  "processing",
					"amount":  100.00 * orderCount,
				})
				evt := event.Event{
					Type:    event.Broadcast,
					Payload: payload,
				}
				_ = m.SendEvent(evt)
				log.Info("发送order消息: order-%d", orderCount)

			case <-tickerNotice.C:
				noticeCount++
				payload, _ := event.New("notice", map[string]any{
					"title":   fmt.Sprintf("通知 #%d", noticeCount),
					"content": "您有一条新消息",
				})
				evt := event.Event{
					Type:    event.Broadcast,
					Payload: payload,
				}
				_ = m.SendEvent(evt)
				log.Info("发送notice消息: notice-%d", noticeCount)

			case <-tickerSystem.C:
				systemCount++
				payload, _ := event.New("system", map[string]any{
					"level":   "info",
					"message": fmt.Sprintf("系统消息 #%d", systemCount),
				})
				evt := event.Event{
					Type:    event.Broadcast,
					Payload: payload,
				}
				_ = m.SendEvent(evt)
				log.Info("发送system消息: system-%d", systemCount)
			}
		}
	}()

	log.Info("startup web server with event support...")
	log.Info("请使用以下URL连接SSE并订阅event:")
	log.Info("  - 订阅order_update: http://localhost:8080/sse?events=order_update")
	log.Info("  - 订阅new_notification: http://localhost:8080/sse?events=new_notification")
	log.Info("  - 订阅多个: http://localhost:8080/sse?events=order_update,new_notification,system_alert")
	_ = http.ListenAndServe(":8080", mux)
}
