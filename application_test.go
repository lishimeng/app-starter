package app

import (
	"testing"
)

func TestNew(t *testing.T) {
	var a Application
	target  := New()
	a = target.(Application)
	t.Logf("application:%T", a)
}

func TestGetOrm(t *testing.T) {
	o := GetOrm()
	t.Logf("Orm is :%T", o)
}

func TestGetCache(t *testing.T) {
	o := GetCache()
	t.Logf("Cache is :%T", o)
}

