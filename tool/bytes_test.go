package tool

import "testing"

func TestBytesToHex(t *testing.T) {
	bs := []byte{0x01, 0x02, 0x04, 0x5F}
	t.Log(BytesToHex(bs, "0x", ":"))
	t.Log(BytesToHex(bs, "0X", ":"))
	t.Log(BytesToHex(bs, "", " "))
	t.Log(BytesToHex(bs, "", ":"))
}
