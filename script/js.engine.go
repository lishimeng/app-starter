package script

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"time"
)

type JsEngine interface {
	Invoke(method string, params ...interface{}) (otto.Value, error)
	Inject(name string, callback func(call otto.FunctionCall) otto.Value)
	SetValue(name string, value interface{}) error
}

func Create(javascript string) (engine JsEngine, err error) {

	return CreateWithTimeout(javascript, 10*time.Millisecond)
}

func CreateWithTimeout(javascript string, timeout time.Duration) (engine JsEngine, err error) {

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)
	vm.SetStackDepthLimit(32)
	err = LoadScript(vm, javascript, timeout)

	if err == nil {
		js := jsEngine{
			script:           javascript,
			vm:               vm,
			maxExecutionTime: timeout,
		}
		engine = &js
	}

	return engine, err
}

func LoadScript(vm *otto.Otto, script string, timeout time.Duration) (err error) {
	defer func() {
		if exp := recover(); exp != nil {
			err = fmt.Errorf("%s", exp)
		}
	}()
	go func() {
		time.Sleep(timeout)
		vm.Interrupt <- func() {
			panic(errors.New("execute javascript timeout"))
		}
	}()
	_, err = vm.Run(script)
	return err
}

func CallFunc(vm *otto.Otto, method string, timeout time.Duration, params ...interface{}) (value otto.Value, err error) {
	defer func() {
		if exp := recover(); exp != nil {
			err = fmt.Errorf("%s", exp)
		}
	}()

	// 执行js前打开线程做超时检查
	go func() {
		time.Sleep(timeout)
		vm.Interrupt <- func() {
			panic(errors.New("execute javascript timeout"))
		}
	}()

	value, err = vm.Call(method, nil, params...)
	return value, err
}

func Execute(vm *otto.Otto, script string, timeout time.Duration) (value otto.Value, err error) {

	defer func() {
		if exp := recover(); exp != nil {
			err = fmt.Errorf("%s", exp)
		}
	}()

	// 执行js前打开线程做超时检查
	go func() {
		time.Sleep(timeout)
		vm.Interrupt <- func() {
			panic(errors.New("execute javascript timeout"))
		}
	}()

	value, err = vm.Run(script)
	if err != nil {
		return value, err
	}
	if value.IsFunction() {
		return value, errors.New("not support javascript return type is 'function'")
	}

	return value, err
}
