package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/server"
	shutdown "github.com/lishimeng/go-app-shutdown"
)

func TestWebGlobalMiddlewareWithRequest(t *testing.T) {
	var logBuf bytes.Buffer
	log.Config().LevelFromString("INFO").Text().Out(&logBuf).Apply()
	t.Cleanup(func() {
		log.Config().LevelFromString("INFO").Text().Out(io.Discard).Apply()
	})

	webAddr := freeTCPAddr(t)
	done := make(chan struct{})

	go func() {
		defer close(done)
		time.Sleep(250 * time.Millisecond)

		resp, err := http.Get("http://" + webAddr + "/hello")
		if err != nil {
			t.Errorf("GET /hello: %v", err)
			shutdown.Exit("request failed")
			return
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, body %s", resp.StatusCode, body)
		}
		if resp.Header.Get("X-Global-Mw") != "1" {
			t.Errorf("missing X-Global-Mw, headers=%v", resp.Header)
		}
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("json: %v body=%s", err, body)
		} else if payload["msg"] != "ok" {
			t.Errorf("payload = %#v", payload)
		}

		logs := logBuf.String()
		t.Logf("captured logs:\n%s", logs)
		if !strings.Contains(logs, "global middleware hit") {
			t.Errorf("log missing middleware message, got:\n%s", logs)
		}
		if !strings.Contains(logs, "/hello") {
			t.Errorf("log missing path /hello, got:\n%s", logs)
		}

		shutdown.Exit("test done")
	}()

	err := New().Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.
			DisableAdmin().
			EnableWebMiddleware(func(c server.Context) {
				log.Info("global middleware hit", "path", c.Path(), "method", c.Request().Method)
				c.SetHeader("X-Global-Mw", "1")
				c.Next()
			}).
			EnableWeb(webAddr, func(root server.Router) {
				root.Get("/hello", func(c server.Context) {
					log.Info("handler hello")
					c.JSON(map[string]string{"msg": "ok"})
				})
			})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	<-done

	if t.Failed() {
		t.FailNow()
	}
}

func freeTCPAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}
