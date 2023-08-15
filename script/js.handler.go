package script

import (
	"github.com/robertkrimen/otto"
	"time"
)

type jsEngine struct {
	script           string
	vm               *otto.Otto
	maxExecutionTime time.Duration
}

func (engine *jsEngine) Invoke(method string, params ...interface{}) (otto.Value, error) {
	res, err := CallFunc(engine.vm, method, engine.maxExecutionTime, params...)
	return res, err
}

func (engine *jsEngine) Inject(name string, callback func(call otto.FunctionCall) otto.Value) {
	_ = engine.vm.Set(name, callback)
}

func (engine *jsEngine) SetValue(name string, value interface{}) error {
	return engine.vm.Set(name, value)
}
