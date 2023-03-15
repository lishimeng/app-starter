package tool

import (
	"context"
	"time"
)

func NewInterval(ctx context.Context, duration time.Duration, exitOnErr bool, task func() error) {
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
				err := task()
				if err != nil && exitOnErr {
					timer.Stop()
				} else {
					timer.Reset(duration)
				}
			}
		}
	}()
}
