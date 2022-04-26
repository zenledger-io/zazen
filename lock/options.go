package lock

import (
	"math"
	"time"
)

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
		o.RetryRandMax = time.Duration(math.Round(float64(o.RetryDurMax) * 0.25))
	}
}
