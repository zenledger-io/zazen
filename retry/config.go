package retry

import (
	"context"
	"github.com/zenledger-io/go-utils/timeutils"
	"math/rand"
	"time"
)

var (
	DefaultBackoffFactor = 0.3
	DefaultRandomFactor  = 0.2
)

type ShouldRetryFunc func(ctx context.Context, i int, err error) bool

type Config struct {
	ShouldRetry   ShouldRetryFunc
	Wait          time.Duration
	MaxWait       time.Duration
	BackOffFactor float64
	RandomFactor  float64
}

func NewConfig(wait time.Duration, maxWait time.Duration, shouldRetry ShouldRetryFunc) Config {
	return Config{
		ShouldRetry:   shouldRetry,
		Wait:          wait,
		MaxWait:       maxWait,
		BackOffFactor: DefaultBackoffFactor,
		RandomFactor:  DefaultRandomFactor,
	}
}

func (cfg Config) NextWait(i int) time.Duration {
	f := float64(cfg.Wait)
	if cfg.BackOffFactor > 0 {
		for j := 0; j < i; j++ {
			f += f * cfg.BackOffFactor
		}
	}

	f += float64(cfg.Wait) * cfg.RandomFactor * (rand.Float64() - 0.5)

	wait := time.Duration(f)
	if cfg.MaxWait > 0 && wait > cfg.MaxWait {
		wait = cfg.MaxWait
	}

	return wait
}

func (cfg Config) Sleep(ctx context.Context, i int) {
	wait := cfg.NextWait(i)
	if wait <= 0 || ctx.Err() != nil {
		return
	}

	t := time.NewTimer(wait)
	defer timeutils.StopTimer(t)

	select {
	case <-ctx.Done():
		return
	case <-t.C:
		return
	}
}

func (cfg Config) Do(ctx context.Context, f func(context.Context) error) (int, error) {
	return Do(ctx, cfg, f)
}
