package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/zenledger-io/go-utils/lock"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRedisLocker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	defer func(dur time.Duration) {
		lock.DefaultRetryRandMax = dur
	}(lock.DefaultRetryRandMax)
	lock.DefaultRetryRandMax = 500 * time.Nanosecond

	radd := os.Getenv("REDIS_ADDR")
	opts := redis.Options{Addr: radd}
	rc := redis.NewClient(&opts)
	defer rc.Close()

	var val int
	lopts := lock.Options{
		MaxRetry:    10,
		RetryDurMin: 1 * time.Millisecond,
		RetryDurMax: 20 * time.Millisecond,
	}
	locker := lock.NewRedisLocker(rc, "testkey", 3*time.Second, &lopts)
	runs := 100
	wg := sync.WaitGroup{}
	wg.Add(runs)
	for i := 0; i < runs; i++ {
		go func() {
			defer wg.Done()

			ulock, err := locker.Lock(ctx)
			require.NoError(t, err)
			defer ulock.Unlock()

			val += 1
		}()
	}
	wg.Wait()

	require.Equal(t, runs, val)
}
