package tool

import (
	"context"
	"time"
)

func NewInterval(ctx context.Context, duration time.Duration, task func()) {
	go func() {
		var timer = time.NewTimer(duration)
		defer func() {
			timer.Stop()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				task()
				timer.Reset(duration)
			}
		}
	}()
}
