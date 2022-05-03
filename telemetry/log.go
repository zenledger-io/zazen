package telemetry

import (
	"context"
)

type ctxKey int

const logCtxKey ctxKey = iota

// Field is a structured logging field..
type Field struct {
	Name  string
	Value any
}

// Log is a interface for structure logging.
type Log interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, err error, fields ...Field)
}

// ContextWithLog returns a new context that carries the supplied Log.
func ContextWithLog(ctx context.Context, l Log) context.Context {
	return context.WithValue(ctx, logCtxKey, l)
}

// LogFromContext returns the Log value stored in ctx, if any.
// If a value is not found, a nop Log is returned.
func LogFromContext(ctx context.Context) Log {
	if log, ok := ctx.Value(logCtxKey).(Log); ok {
		return log
	}
	return nopLog{}
}

type nopLog struct{}

func (nopLog) Debug(context.Context, string, ...Field)        {}
func (nopLog) Info(context.Context, string, ...Field)         {}
func (nopLog) Error(context.Context, string, error, ...Field) {}
