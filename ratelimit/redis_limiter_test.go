package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/zenledger-io/zazen/internal/testing/config"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRedisLimiter_Allow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	radd := os.Getenv("REDIS_ADDR")
	opts := redis.Options{Addr: radd, DB: config.RatelimitRedisDB}
	rc := redis.NewClient(&opts)
	defer rc.Close()

	rate := 10
	period := 500 * time.Millisecond
	var val int64
	limiter := NewRedisLimiter(rc, "TestRedisLimiter_Allow", rate, period)
	runs := rate + 1
	wg := sync.WaitGroup{}
	wg.Add(runs)
	start := time.Now()
	for i := 0; i < runs; i++ {
		go func() {
			defer wg.Done()

			retryAfter, err := limiter.Allow(ctx)
			require.NoError(t, err)

			for retryAfter > 0 {
				time.Sleep(retryAfter)
				retryAfter, err = limiter.Allow(ctx)
				require.NoError(t, err)
			}

			atomic.AddInt64(&val, 1)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	require.GreaterOrEqual(t, elapsed, period)
	require.Equal(t, runs, int(val))
}
