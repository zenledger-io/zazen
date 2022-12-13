package metrics

import (
	"context"
)

type Transaction interface {
	End()
	AddAttribute(key string, value any)
	AddAttributes(attrs map[string]any)
	StartTransaction(ctx context.Context, name string) (Transaction, context.Context)
	NoticeError(err error)
}

func StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	tx, ok := ctx.Value(txKey).(Transaction)
	if !ok {
		return ContextMonitor(ctx).StartTransaction(ctx, name)
	}

	return tx.StartTransaction(ctx, name)
}
