package mqtt

import (
	"fmt"
	"github.com/lishimeng/x/util"
	"strings"
)

type BuilderOption struct {
	Share bool   // 是否共享topic
	Tpl   bool   // true: /a/b/{c}/d  false:/a/b/c/d
	Group string // group名称, 默认default
}

const (
	shareTopicPrefix = "$share/"
	TopicGroup       = "default"
)

func TopicResolver(tpl string, topic string) (res map[string]string, err error) {

	ss := strings.Split(tpl, "/")
	st := strings.Split(topic, "/")
	if len(ss) == 0 || len(st) != len(ss) {
		err = fmt.Errorf("topics is not match the template %s[%s]", tpl, topic)
	} else {
		res = make(map[string]string)
		for i, v := range ss {
			if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
				name := v[1 : len(v)-1]
				value := st[i]
				res[name] = value
			}
		}
	}
	return res, err
}

func TopicBuilder(opt BuilderOption, format string, key ...any) (t string) {
	var tmp []any
	var group = TopicGroup
	if len(opt.Group) > 0 {
		group = opt.Group
	}
	if opt.Tpl {
		for _, k := range key {
			t = fmt.Sprintf("{%s}", k)
			tmp = append(tmp, t)
		}
	} else {
		tmp = key
	}

	t = fmt.Sprintf(format, tmp...)
	if opt.Share {
		if strings.HasPrefix(t, "/") {
			t = util.Join("", shareTopicPrefix, group, t)
		} else {
			t = util.Join("", shareTopicPrefix, group, "/", t)
		}
	}
	return
}
