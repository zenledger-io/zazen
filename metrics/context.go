package metrics

import "context"

var (
	monitorKey contextKey = "monitor"
	txKey      contextKey = "transaction"
)

type contextKey string

func ContextWithMonitor(ctx context.Context, m Monitor) context.Context {
	return context.WithValue(ctx, monitorKey, m)
}

func ContextMonitor(ctx context.Context) Monitor {
	if m, ok := ctx.Value(monitorKey).(Monitor); ok {
		return m
	}

	return NewNullMonitor()
}

func ContextWithTransaction(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ContextTransaction(ctx context.Context) Transaction {
	if tx, ok := ctx.Value(txKey).(Transaction); ok {
		return tx
	}

	return nullTransaction{}
}
