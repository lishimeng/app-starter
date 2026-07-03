package api

import (
	"net/http"
	"testing"
)

func TestDefaultLiveness(t *testing.T) {
	if code := DefaultLiveness(); code != http.StatusOK {
		t.Fatalf("code = %d, want 200", code)
	}
}

func TestDefaultReadinessNoDeps(t *testing.T) {
	if code := DefaultReadiness(); code != http.StatusOK {
		t.Fatalf("code = %d, want 200 without deps", code)
	}
}
