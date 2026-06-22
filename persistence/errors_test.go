package persistence

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestNormalizeErr_NotFound(t *testing.T) {
	got := NormalizeErr(gorm.ErrRecordNotFound)
	if !IsNotFound(got) {
		t.Fatalf("expected not found, got %v", got)
	}
}

func TestNormalizeErr_Nil(t *testing.T) {
	if NormalizeErr(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestNormalizeErr_Other(t *testing.T) {
	raw := errors.New("boom")
	if NormalizeErr(raw) != raw {
		t.Fatal("expected same error")
	}
}

func TestIsDuplicate(t *testing.T) {
	if !IsDuplicate(errors.New("duplicate key value violates unique constraint")) {
		t.Fatal("expected duplicate")
	}
	if IsDuplicate(errors.New("other")) {
		t.Fatal("expected not duplicate")
	}
}
