package lock

import (
	"context"
	"math"
	"math/rand"
	"time"

	lock "github.com/bsm/redislock"
)

var (
	DefaultRetryRandMax = 1 * time.Millisecond
)

func NewRedisLocker(c lock.RedisClient, key string, ttl time.Duration, opts *Options) Locker {
	opts.SetDefaults()

	return &redisLocker{
		lc:   lock.New(c),
		key:  key,
		opts: opts,
		ttl:  ttl,
	}
}

type redisLocker struct {
	lc   *lock.Client
	key  string
	opts *Options
	ttl  time.Duration
}

func (l *redisLocker) Lock(ctx context.Context) (Unlocker, error) {
	var lockOpts *lock.Options
	if l.opts != nil {
		strategy := lock.ExponentialBackoff(l.opts.RetryDurMin, l.opts.RetryDurMax)
		strategy = lock.LimitRetry(strategy, l.opts.MaxRetry)
		lockOpts = &lock.Options{
			RetryStrategy: &randomizedRetry{
				MaxRand:       l.opts.RetryRandMax,
				RetryStrategy: strategy,
			},
		}
	}
	rl, err := l.lc.Obtain(ctx, l.key, l.ttl, lockOpts)
	if err != nil {
		return nil, err
	}

	return redisUnlocker{rl}, nil
}

type redisUnlocker struct {
	lock *lock.Lock
}

func (l redisUnlocker) Unlock() error {
	return l.lock.Release(context.Background())
}

type randomizedRetry struct {
	MaxRand time.Duration
	lock.RetryStrategy
}

func (rr *randomizedRetry) NextBackoff() time.Duration {
	dur := rr.RetryStrategy.NextBackoff()
	rFactor := float64(rr.MaxRand) * rand.Float64()
	dur += time.Duration(math.Round(rFactor))
	return dur
}
