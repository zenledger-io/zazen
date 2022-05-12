package lock

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	lock "github.com/bsm/redislock"
)

var (
	RedisLockNamespace = "_lock"
)

func NewRedisLocker(c lock.RedisClient, key string, ttl time.Duration, opts *Options) Locker {
	return newRedisLocker(c, key, ttl, opts)
}

func NewRedisRefreshLocker(c lock.RedisClient, key string, ttl time.Duration, opts *Options) RefreshLocker {
	return newRedisLocker(c, key, ttl, opts)
}

func newRedisLocker(c lock.RedisClient, key string, ttl time.Duration, opts *Options) *redisLocker {
	return &redisLocker{
		lc:   lock.New(c),
		key:  fmt.Sprintf("%v:%v", RedisLockNamespace, key),
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
	return l.lock(ctx)
}

func (l *redisLocker) LockWithRefresh(ctx context.Context) (RefreshUnlocker, error) {
	return l.lock(ctx)
}

func (l *redisLocker) Indefinite(ctx context.Context, f func(context.Context)) error {
	return Indefinite(ctx, l, l.ttl*3/5, f)
}

func (l *redisLocker) lock(ctx context.Context) (*redisUnlocker, error) {
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

	ul := l.unlocker(rl, lockOpts)
	return &ul, nil
}

func (l *redisLocker) unlocker(rl *lock.Lock, opts *lock.Options) redisUnlocker {
	return redisUnlocker{
		unlock: func() error {
			return rl.Release(context.Background())
		},
		refresh: func(ctx context.Context) error {
			return rl.Refresh(ctx, l.ttl, opts)
		},
	}
}

type redisUnlocker struct {
	unlock  func() error
	refresh func(context.Context) error
}

func (l redisUnlocker) Unlock() error {
	return l.unlock()
}

func (l redisUnlocker) Refresh(ctx context.Context) error {
	return l.refresh(ctx)
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
