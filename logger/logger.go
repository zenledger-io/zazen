package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const loggerCtxKey ctxKey = 0

// NewContext returns a new Context that carries a value l.
func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// FromContext returns the *zap.Logger value stored in ctx, if any.
// If a value is not found, a nop logger is returned.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerCtxKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
