package retry

import (
	"context"
	"time"
)

type ShouldRetryFunc func(ctx context.Context, i int, err error) bool

func (srf ShouldRetryFunc) Wrap(f ShouldRetryFunc) ShouldRetryFunc {
	if srf == nil {
		return f
	}

	return func(ctx context.Context, i int, err error) bool {
		return srf(ctx, i, err) && f(ctx, i, err)
	}
}

func WrapShouldRetryFuncs(fs ...ShouldRetryFunc) ShouldRetryFunc {
	var f ShouldRetryFunc
	for _, srf := range fs {
		f = f.Wrap(srf)
	}

	return f
}

func MaxRetries(max int) ShouldRetryFunc {
	return func(ctx context.Context, i int, err error) bool {
		if i < max {
			return true
		}

		return false
	}
}

type ResetCountFunc func(int, time.Duration) int

func ResetCountAfter(dur time.Duration) ResetCountFunc {
	return func(i int, sinceLast time.Duration) int {
		if sinceLast >= dur {
			return 0
		}

		return i
	}
}
