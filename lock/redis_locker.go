package lock

import (
	"context"
	"math"
	"math/rand"
	"time"

	lock "github.com/bsm/redislock"
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
	lock.RetryStrategy
}

func (rr *randomizedRetry) NextBackoff() time.Duration {
	dur := rr.RetryStrategy.NextBackoff()
	if dur <= 0 {
		return dur
	}

	rFactor := float64(dur) * rand.Float64() * 0.5
	dur -= time.Duration(math.Round(rFactor))
	return dur
}
