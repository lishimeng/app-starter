package tool

import "testing"

func TestJoin(t *testing.T) {
	delimiter := "/"
	t.Log(Join(delimiter, "sss", "bbb", "mmm"))
}
