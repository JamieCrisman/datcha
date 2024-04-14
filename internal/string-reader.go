package internal

import (
	"context"
	"io"
	"time"
)

func StringRepeaterReader(ctx context.Context, d time.Duration, s string) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		t := time.NewTicker(d)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				pw.Close()
				return
			case <-t.C:
				if _, err := io.WriteString(pw, s); err != nil {
					return
				}
			}
		}
	}()
	return pr
}
