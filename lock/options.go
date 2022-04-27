package lock

import (
	"time"
)

type Options struct {
	MaxRetry    int
	RetryDurMin time.Duration
	RetryDurMax time.Duration
}
