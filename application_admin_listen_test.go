package app

import (
	"testing"

	"github.com/lishimeng/app-starter/server"
)

func TestAdminListenAddr_Default(t *testing.T) {
	b := &ApplicationBuilder{}
	if got := b.adminListenAddr(); got != server.DefaultAdminListen {
		t.Fatalf("no SetAdminListen: got %q, want %q", got, server.DefaultAdminListen)
	}
}

func TestAdminListenAddr_EmptyArgUsesDefault(t *testing.T) {
	b := &ApplicationBuilder{}
	b.SetAdminListen("")
	if got := b.adminListenAddr(); got != server.DefaultAdminListen {
		t.Fatalf("SetAdminListen(\"\"): got %q, want %q", got, server.DefaultAdminListen)
	}
}

func TestAdminListenAddr_Custom(t *testing.T) {
	b := &ApplicationBuilder{}
	b.SetAdminListen(":7070")
	if got := b.adminListenAddr(); got != ":7070" {
		t.Fatalf("SetAdminListen(\":7070\"): got %q, want :7070", got)
	}
}

func TestAdminListenAddr_DisableAdmin(t *testing.T) {
	t.Run("disable only", func(t *testing.T) {
		b := &ApplicationBuilder{}
		b.DisableAdmin()
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("DisableAdmin: got %q, want empty", got)
		}
	})
	t.Run("set then disable", func(t *testing.T) {
		b := &ApplicationBuilder{}
		b.SetAdminListen(":7070")
		b.DisableAdmin()
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("SetAdminListen then DisableAdmin: got %q, want empty", got)
		}
	})
	t.Run("set empty then disable", func(t *testing.T) {
		b := &ApplicationBuilder{}
		b.SetAdminListen("")
		b.DisableAdmin()
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("SetAdminListen(\"\") then DisableAdmin: got %q, want empty", got)
		}
	})
}
