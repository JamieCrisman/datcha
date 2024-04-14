package tick

import (
	"context"
	"time"
)

func Do(ctx context.Context, d time.Duration, f func()) {
	go func() {
		t := time.NewTicker(d)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				f()
			}
		}
	}()
}
