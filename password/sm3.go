package passwd

import (
	"github.com/ZZMarquis/gm/sm3"
)

func gen(plain []byte) (bs []byte, err error) {
	p := sm3.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	bs = p.Sum(nil)
	return
}
