package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContextQueryAndParam(t *testing.T) {
	srv := New(Config{Listen: ":0"})
	srv.RegisterComponent(func(root Router) {
		root.Get("/items/:id", func(ctx Context) {
			id, err := ctx.ParamInt("id")
			if err != nil {
				ctx.Status(http.StatusBadRequest)
				return
			}
			page, _ := ctx.QueryInt("page")
			ctx.Json(map[string]any{"id": id, "page": page})
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/items/42?page=3", nil)
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["id"].(float64) != 42 {
		t.Fatalf("id = %v, want 42", body["id"])
	}
	if body["page"].(float64) != 3 {
		t.Fatalf("page = %v, want 3", body["page"])
	}
}

func TestContextMiddlewareSetGet(t *testing.T) {
	srv := New(Config{Listen: ":0"})
	srv.RegisterComponent(func(root Router) {
		root.Get("/mw", func(ctx Context) {
			ctx.Set("key", "value")
			ctx.Next()
		}, func(ctx Context) {
			v, ok := ctx.Get("key")
			if !ok || v.(string) != "value" {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			ctx.Json(map[string]string{"ok": "true"})
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/mw", nil)
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestContextBindJSON(t *testing.T) {
	srv := New(Config{Listen: ":0"})
	srv.RegisterComponent(func(root Router) {
		root.Post("/echo", func(ctx Context) {
			var payload struct {
				Name string `json:"name"`
			}
			if err := ctx.BindJSON(&payload); err != nil {
				ctx.Status(http.StatusBadRequest)
				return
			}
			ctx.Json(payload)
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.GetEngine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Name != "test" {
		t.Fatalf("name = %q, want test", body.Name)
	}
}
