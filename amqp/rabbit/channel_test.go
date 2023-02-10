package rabbit

import (
	"context"
	"github.com/lishimeng/go-log"
	"testing"
	"time"
)

func TestChannel001(t *testing.T) {

	var ctx = context.TODO()

	var c = make(chan byte, 100)

	go func() { // read_1
		for {
			select {
			case <-ctx.Done():
				log.Info("read_1 done")
				return
			case v := <-c:
				log.Info("read_1 <- %+v\n", v)
			}
		}
	}()

	go func() { // read_2
		for {
			select {
			case <-ctx.Done():
				log.Info("read_2 done")
				return
			case v := <-c:
				log.Info("read_2 <- %+v\n", v)
			}
		}
	}()

	var index byte = 0x00

	for i := 0; i < 300; i++ {
		c <- index
		index = index + 1
	}

	time.Sleep(time.Second * 10)
	ctx.Done()
}
