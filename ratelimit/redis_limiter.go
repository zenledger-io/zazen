package ratelimit

import (
	"context"
	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"time"
)

// Rediser is a copied from a non exported interface in github.com/go-redis/redis_rate/v10
type Rediser interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
}

func NewRedisLimiter(r Rediser, key string, rate int, period time.Duration) Limiter {
	return &redisLimiter{
		limiter: redisrate.NewLimiter(r),
		key:     key,
		rate:    rate,
		period:  period,
	}
}

type redisLimiter struct {
	limiter *redisrate.Limiter
	key     string
	rate    int
	period  time.Duration
}

func (l *redisLimiter) Allow(ctx context.Context) (time.Duration, error) {
	r, err := l.limiter.Allow(ctx, l.key, redisrate.Limit{
		Rate:   l.rate,
		Burst:  1,
		Period: l.period,
	})
	if err != nil {
		return 0, err
	}

	return r.RetryAfter, nil
}
