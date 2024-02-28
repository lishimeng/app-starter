package stream

import (
	"bytes"
	"context"
	"github.com/lishimeng/x/util"
	"testing"
	"time"
)

func TestSplitPP(t *testing.T) {
	var str = "dfasdasfa;fasdf;df;asd;fdfasdf;df;fddsg34" //最后一段不处理
	spp := NewSplitPP([]byte(";")[0])
	spp.Listen(func(p []byte) {
		t.Log(string(p))
	})
	n := spp.Data([]byte(str))
	t.Log(n)
}

func TestHeadTail(t *testing.T) {
	//var str = "[[5555]][[2222]][[4444444]][[22222222" //最后一段不处理
	var str = "[5][2][4][22222222" //最后一段不处理
	spp := NewHeadTailPP([]byte("["), []byte("]"))
	spp.Listen(func(p []byte) {
		t.Log(string(p))
	})
	n := spp.Data([]byte(str))
	t.Log(n)
}

func TestSerial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	pp := NewHeadTailPP([]byte{0x68}, []byte{0x16})
	pp.Listen(func(p []byte) {
		t.Log(util.BytesToHex(p, "-"))
	})
	s, err := NewSerialSession(ctx, "COM5", 115200, WithReact(), WithPacketProcessor(pp))

	if err != nil {
		t.Fatal(err)
	}
	_, _ = s.Write([]byte{0x68, 0x02, 0x05, 0x0A, 0x00, 0x00, 0x00, 0x00, 0xB2, 0x16})
	time.Sleep(time.Second * 2)
	cancel()
	time.Sleep(time.Second * 1)
	t.Log("done")
}

func TestBuff(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte("fsadfdfdfa"))
	var p = make([]byte, 100)
	n, err := buf.Read(p)
	//n, err = buf.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(p[:n]))
	t.Log(n)
}
