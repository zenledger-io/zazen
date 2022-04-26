package lock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(context.Context) (Unlocker, error)
}

type Unlocker interface {
	Unlock() error
}

type Options struct {
	MaxRetry     int
	RetryDurMin  time.Duration
	RetryDurMax  time.Duration
	RetryRandMax time.Duration
}

func (o *Options) SetDefaults() {
	if o == nil {
		return
	}

	if o.RetryRandMax == 0 {
		o.RetryRandMax = DefaultRetryRandMax
	}
}
