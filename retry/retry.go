package retry

import (
	"context"
	"time"
)

func Do(ctx context.Context, cfg Config, f func(context.Context) error) (int, error) {
	var lastErrorAt time.Time
	var i int
	for {
		if err := ctx.Err(); err != nil {
			return i, err
		}

		if err := f(ctx); err != nil {
			if cfg.ResetCount != nil {
				i = cfg.ResetCount(i, time.Since(lastErrorAt))
				lastErrorAt = time.Now()
			}

			if cfg.ShouldRetry == nil || cfg.ShouldRetry(ctx, i, err) {
				i += 1
				cfg.Sleep(ctx, i)
				continue
			}

			return i, err
		} else {
			return i, nil
		}
	}
}
