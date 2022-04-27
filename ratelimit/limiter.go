package ratelimit

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context) (time.Duration, error)
}
