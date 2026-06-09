package persistence

import "testing"

func TestCondStrSkipsEmpty(t *testing.T) {
	m := &mockQuery{}
	m.EqualStr("code", "")
	m.LikeStr("name", "")
	m.ILikeStr("name", "")
	if m.filters != 0 {
		t.Fatalf("expected 0 filters, got %d", m.filters)
	}
	m.EqualStr("code", "abc").LikeStr("name", "x")
	if m.filters != 2 {
		t.Fatalf("expected 2 filters, got %d", m.filters)
	}
}

func TestCondBuildersChain(t *testing.T) {
	m := &mockQuery{}
	m.Equal("status", 1).
		NotEqual("deleted", 1).
		In("id", []int{1, 2}).
		Like("name", "a").
		LLike("name", "a").
		RLike("name", "a").
		ILike("name", "a")
	if m.filters != 7 {
		t.Fatalf("expected 7 filters, got %d", m.filters)
	}
}
