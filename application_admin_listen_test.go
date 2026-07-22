package app

import (
	"testing"

	"github.com/lishimeng/app-starter/server"
)

func enabledAdminBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{adminEnable: true}
}

func TestAdminListenAddr_Default(t *testing.T) {
	b := enabledAdminBuilder()
	if got := b.adminListenAddr(); got != server.DefaultAdminListen {
		t.Fatalf("no SetAdminListen: got %q, want %q", got, server.DefaultAdminListen)
	}
}

func TestAdminListenAddr_EmptyArgUsesDefault(t *testing.T) {
	b := enabledAdminBuilder()
	b.SetAdminListen("")
	if got := b.adminListenAddr(); got != server.DefaultAdminListen {
		t.Fatalf("SetAdminListen(\"\"): got %q, want %q", got, server.DefaultAdminListen)
	}
}

func TestAdminListenAddr_Custom(t *testing.T) {
	b := enabledAdminBuilder()
	b.SetAdminListen(":7070")
	if got := b.adminListenAddr(); got != ":7070" {
		t.Fatalf("SetAdminListen(\":7070\"): got %q, want :7070", got)
	}
}

func TestAdminListenAddr_DisableAdmin(t *testing.T) {
	t.Run("disable only", func(t *testing.T) {
		b := enabledAdminBuilder()
		b.DisableAdmin()
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("DisableAdmin: got %q, want empty", got)
		}
	})
	t.Run("set then disable", func(t *testing.T) {
		b := enabledAdminBuilder()
		b.SetAdminListen(":7070")
		b.DisableAdmin()
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("SetAdminListen then DisableAdmin: got %q, want empty", got)
		}
	})
	t.Run("disable then set still off", func(t *testing.T) {
		b := enabledAdminBuilder()
		b.DisableAdmin()
		b.SetAdminListen(":7070")
		if got := b.adminListenAddr(); got != "" {
			t.Fatalf("DisableAdmin then SetAdminListen must stay off, got %q", got)
		}
		if b.adminListen != ":7070" {
			t.Fatalf("listen should still be stored as :7070, got %q", b.adminListen)
		}
	})
	t.Run("disable then enable restores", func(t *testing.T) {
		b := enabledAdminBuilder()
		b.SetAdminListen(":7070")
		b.DisableAdmin()
		b.EnableAdmin()
		if got := b.adminListenAddr(); got != ":7070" {
			t.Fatalf("EnableAdmin after disable: got %q, want :7070", got)
		}
	})
}

func TestAdminListenAddr_DisabledByDefaultZeroValue(t *testing.T) {
	b := &ApplicationBuilder{}
	if got := b.adminListenAddr(); got != "" {
		t.Fatalf("zero-value builder adminEnable=false: got %q, want empty", got)
	}
}
