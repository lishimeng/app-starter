package script

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"testing"
	"time"
)

func TestJsEngine_Invoke(t *testing.T) {
	testContent := `function decode(fport, data) {
return {"a": 12, "b": "ffdasf"}
	}`

	var vm, err = Create(testContent)
	if err != nil {
		return
	}
	raw, err := vm.Invoke("decode", "", "")
	if err != nil {
		return
	}
	ras, err := raw.Export()
	if err != nil {
		return
	}
	switch r := ras.(type) {
	case map[string]interface{}:
		result := r
		fmt.Println(result)
	default:
		fmt.Println("type err")
	}
}

func TestExecute(t *testing.T) {
	testContent := `function decode(fport, data) {
return {"a": 12, "b": "ffdasf"}
	}`
	vm := otto.New()
	raw, err := Execute(vm, testContent, time.Second)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(raw)
}
