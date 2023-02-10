package rabbit

import (
	"math/rand"
	"time"
)

type Delay struct {
	min             int
	max             int
	cur             int
	randomIncrement bool
}

func (delay *Delay) Reset() {
	delay.cur = 1
}

func (delay *Delay) Delay(after ...func(curDelayDuration int)) {
	select {
	case <-time.After(time.Duration(delay.cur) * time.Second):
		if delay.randomIncrement {
			delay.cur = delay.cur + rand.Intn(5)
		} else {
			delay.cur = delay.cur + 1
		}
		if delay.cur > delay.max {
			delay.cur = delay.max
		}
		if len(after) > 0 {
			after[0](delay.cur)
		}
	}
}

func NewDelay(min int, max int, randomIncrement bool) (d *Delay) {
	d = &Delay{
		min:             min,
		max:             max,
		randomIncrement: randomIncrement,
	}
	return
}
