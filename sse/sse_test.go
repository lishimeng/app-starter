package sse

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"testing"
	"time"

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
		WebHandler(m, w, r)
	})
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Info("提取static目录失败:", err)
	}

	pages := http.FS(staticFS)

	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		// 如果访问根路径，直接返回index.html
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
		for {
			select {
			case <-time.After(time.Second * 10):
				i++
				log.Info("server[broadcast]--> all client %d", i)
				m.Broadcast <- fmt.Sprintf("test msg: %d", i)
			}
		}
	}()

	log.Info("startup web server...")
	_ = http.ListenAndServe(":8080", mux)
}
