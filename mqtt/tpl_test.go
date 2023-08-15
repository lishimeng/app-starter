package mqtt

import (
	"errors"
	"testing"
)

func TestTplResolve(t *testing.T) {
	var tpl = "aaa/bbb/{para}/d"
	var topic = "aaa/bbb/cccc/d"
	m, err := TopicResolver(tpl, topic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m)
}

func TestTplBuilder(t *testing.T) {
	var format = "aaa/bbb/%s/d/%s/m"
	var key = "Device"
	var key2 = "gateway"
	var tpl = TopicBuilder(BuilderOption{
		Share: false,
		Tpl:   true,
	}, format, key, key2)
	t.Log(tpl)
	if tpl != "aaa/bbb/{Device}/d/{gateway}/m" {
		t.Fatal(errors.New(tpl))
	}
}
