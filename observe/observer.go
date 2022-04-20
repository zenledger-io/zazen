package observe

import (
	"context"
	"fmt"

	"github.com/honeybadger-io/honeybadger-go"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ctxKey int

const observerCtxKey ctxKey = 0

// Field is metadata for an observable event.
type Field struct {
	Name  string
	Value any
}

// Observer is a type that is used to make
// a service observable.
type Observer interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, err error, fields ...Field)
	Panic(ctx context.Context, message string, stack []byte, fields ...Field)
}

// NewContext returns a new Context that carries a value obs.
func NewContext(ctx context.Context, obs Observer) context.Context {
	return context.WithValue(ctx, observerCtxKey, obs)
}

// FromContext returns the Observer value stored in ctx, if any.
// If a value is not found, a nop Observer is returned.
func FromContext(ctx context.Context) Observer {
	if obs, ok := ctx.Value(observerCtxKey).(Observer); ok {
		return obs
	}
	return NewNop()
}

// New creates a new observer.
func New(lgr *zap.Logger, hc *honeybadger.Client) Observer {
	return &observer{
		lgr: lgr,
		hc:  hc,
	}
}

type observer struct {
	lgr *zap.Logger
	hc  *honeybadger.Client
}

func (o *observer) Debug(ctx context.Context, message string, fields ...Field) {
	o.lgr.Debug(message, zapFields(fields)...)
}

func (o *observer) Error(ctx context.Context, message string, err error, fields ...Field) {
	zfs := append(zapFields(fields), zap.Error(err))
	o.lgr.Error(message, zfs...)

	honeybadger.Notify(err, message)

	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.Finish(tracer.WithError(err))
	}
}

func (o *observer) Panic(ctx context.Context, message string, stack []byte, fields ...Field) {
	zfs := append(zapFields(fields), zap.String("stack", string(stack)))
	o.lgr.Panic(message, zfs...)

	honeybadger.Notify(message, string(stack))

	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.Finish(tracer.WithError(
			fmt.Errorf("panic: %s\n%s", message, string(stack))))
	}
}

func zapFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zf = append(zf, zap.Any(f.Name, f.Value))
	}
	return zf
}

// NewNop returns an observer that does nothing.
func NewNop() Observer {
	return &nopObserver{}
}

type nopObserver struct{}

func (*nopObserver) Debug(context.Context, string, ...Field)         {}
func (*nopObserver) Error(context.Context, string, error, ...Field)  {}
func (*nopObserver) Panic(context.Context, string, []byte, ...Field) {}
