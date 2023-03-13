package tool

import "testing"

func TestJoin(t *testing.T) {
	delimiter := "/"
	t.Log(Join(delimiter, "sss", "bbb", "mmm"))
}

func TestGetRandomString(t *testing.T) {
	s := GetRandomString(4)
	t.Log(s)
}
