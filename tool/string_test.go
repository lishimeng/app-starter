package tool

import "testing"

func TestJoin(t *testing.T) {
	delimiter := "/"
	t.Log(Join(delimiter, "sss", "bbb", "mmm"))
}

func TestGetRandHex(t *testing.T) {
	s := RandHexStr(16)
	t.Log(s)
}

func TestGetUUIDString(t *testing.T) {
	t.Log(UUIDString())
}

func TestGetRandStr(t *testing.T) {
	t.Log(RandStr(15))
}
