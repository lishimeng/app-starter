package app

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	var a Application
	ctx := context.Background()
	target  := New(ctx)
	a = target.(Application)
	t.Logf("application:%T", a)
}

func TestGetOrm(t *testing.T) {
	o := GetOrm()
	t.Logf("Orm is :%T", o)
}
