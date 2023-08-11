package factory

import (
	"errors"
	"testing"
)

type testStruct struct {
	Value int
}

type testInter interface {
	GetValue() int
}

func (ts *testStruct) GetValue() int {
	return ts.Value
}

func TestContainerStruct(t *testing.T) {
	var a = testStruct{Value: 15}
	Add(&a, "p")

	var c testStruct
	err := Get(&c, "p")
	if err != nil {
		t.Fatal("not ok, hia hia")
	}
	t.Log(c)
}

func TestContainerInterface(t *testing.T) {
	var b testInter = &testStruct{Value: 15}
	Add(&b)

	var c = new(testInter)
	err := Get(c)
	if err != nil {
		t.Log(err)
		t.Fatal("404, hia hia")
	}
	t.Log((*c).GetValue())
}

func TestContainerLiteral(t *testing.T) {
	var b = 15
	Add(&b, "p")

	var c int
	err := Get(&c, "p")
	if err != nil {
		t.Fatal("not ok, hia hia")
	}
	t.Log(c)
}

func TestPointer001(t *testing.T) {
	var b testInter = &testStruct{Value: 15}
	Add(&b)

	var ptr = new(testInter) // 创建一个数据,而不是定义一个指针变量
	err := Get(ptr)
	if err != nil {
		t.Log(err)
		t.Fatal("404, hia hia")
	}
	t.Log((*ptr).GetValue())
}

func TestPointer002(t *testing.T) {
	var b testInter = &testStruct{Value: 15}
	Add(&b)

	var obj testInter // 创建一个数据
	err := Get(&obj)
	if err != nil {
		t.Log(err)
		t.Fatal("404, hia hia")
	}
	t.Log(obj.GetValue())
}

func TestPointer003(t *testing.T) {
	var b testInter = &testStruct{Value: 15}
	Add(&b)

	var ptr *testInter // 创建一个数据,而不是定义一个指针变量
	err := Get(ptr)
	if errors.Is(err, ErrNotFound) {
		t.Log("yes, result is 404")
		return
	}
	t.Fatal("test fail, expect an err")
}
