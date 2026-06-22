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
	if IsGormRecordNotFound(got) {
		t.Fatal("normalized err must not be gorm.ErrRecordNotFound")
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

func TestIsGormRecordNotFound(t *testing.T) {
	if !IsGormRecordNotFound(gorm.ErrRecordNotFound) {
		t.Fatal("expected gorm.ErrRecordNotFound")
	}
	if IsGormRecordNotFound(NormalizeErr(gorm.ErrRecordNotFound)) {
		t.Fatal("normalized err must not match gorm.ErrRecordNotFound")
	}
	if IsGormRecordNotFound(nil) || IsGormRecordNotFound(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
	}
}

func TestIsNotFound(t *testing.T) {
	normalized := NormalizeErr(gorm.ErrRecordNotFound)
	if !IsNotFound(normalized) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFound(gorm.ErrRecordNotFound) {
		t.Fatal("raw gorm.ErrRecordNotFound is not persistence-level ErrNotFound")
	}
	if !errors.Is(normalized, ErrNotFound) {
		t.Fatal("normalized err must be ErrNotFound")
	}
}

func TestIsNotFoundAny(t *testing.T) {
	if !IsNotFoundAny(gorm.ErrRecordNotFound) {
		t.Fatal("expected raw gorm.ErrRecordNotFound")
	}
	if !IsNotFoundAny(NormalizeErr(gorm.ErrRecordNotFound)) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFoundAny(nil) || IsNotFoundAny(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
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
