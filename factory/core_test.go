package factory

import "testing"

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
	Add(&b, "p")

	var c testInter
	err := Get(&c, "p")
	if err != nil {
		t.Fatal("not ok, hia hia")
	}
	t.Log(c.GetValue())
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
