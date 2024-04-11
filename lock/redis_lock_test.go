package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/zenledger-io/zazen/internal/testing/config"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRedisLocker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	radd := os.Getenv("REDIS_ADDR")
	opts := redis.Options{Addr: radd, DB: config.LockRedisDB}
	rc := redis.NewClient(&opts)
	defer rc.Close()

	var val int
	lopts := Options{
		MaxRetry:    10,
		RetryDurMin: 10 * time.Millisecond,
		RetryDurMax: 50 * time.Millisecond,
	}
	locker := NewRedisLocker(rc, "TestRedisLocker", 3*time.Second, &lopts)
	runs := 100
	wg := sync.WaitGroup{}
	wg.Add(runs)
	for i := 0; i < runs; i++ {
		go func() {
			defer wg.Done()

			ulock, err := locker.Lock(ctx)
			require.NoError(t, err)
			defer ulock.Unlock()

			val += 1 // would fail without lock since -race is used for testing
		}()
	}
	wg.Wait()

	require.Equal(t, runs, val)
}
