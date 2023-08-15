package factory

import (
	"errors"
	"github.com/lishimeng/go-log"
	"reflect"
	"sync"
)

type container struct {
	m     map[string]map[string]interface{}
	ready sync.Mutex
}

const (
	unique = "1"
)

var (
	ErrNotFound = errors.New("404")
)

var c container

func init() {
	c = container{
		m: make(map[string]map[string]interface{}),
	}
}

func Get[T any](ptrType *T, name ...string) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = ErrNotFound
		}
	}()

	c.ready.Lock()
	defer func() {
		c.ready.Unlock()
	}()

	typeName := getTypeName(ptrType)
	var m, ok = c.m[typeName]
	if !ok {
		err = ErrNotFound
		return
	}

	var id string
	if len(name) > 0 {
		id = name[0]
	} else {
		id = unique
	}
	obj, has := m[id]
	if has {
		*ptrType = *obj.(*T) // 刷新到原始值，不是指针变量
		return
	} else {
		err = ErrNotFound
		return
	}
}

func Add[T any](o *T, name ...string) {

	defer func() {
		if e := recover(); e != nil {
			log.Debug(e)
		}
	}()

	c.ready.Lock()
	defer func() {
		c.ready.Unlock()
	}()

	typeName := getTypeName(o)
	var m, ok = c.m[typeName]
	if !ok {
		m = make(map[string]interface{})
		c.m[typeName] = m
	}
	var id string
	if len(name) > 0 {
		id = name[0]
	} else {
		id = unique
	}
	m[id] = o
	return
}

func getTypeName[T any](ptr *T) (name string) {
	t := reflect.TypeOf(ptr)
	return t.Elem().Name()
}
