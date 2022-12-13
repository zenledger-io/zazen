package metrics

import (
	"context"
	"github.com/DataDog/datadog-go/statsd"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
)

var (
	MonitorType      string
	DatadogDebugMode bool
	Addr             string
)

type datadogMonitor struct{}

func NewDatadogMonitor() Monitor {
	return &datadogMonitor{}
}

// Start starts the datadog client and tracer. Make sure the
// DD_AGENT_HOST env variable is set.
func (m *datadogMonitor) Start(ctx context.Context) error {
	client, err := statsd.New(Addr)
	if err != nil {
		return err
	}

	tracer.Start(tracer.WithDebugMode(DatadogDebugMode))

	go func() {
		<-ctx.Done()
		_ = client.Close()
	}()

	return nil
}

func (m *datadogMonitor) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return NewDatadogTransaction(ctx, name)
}

func (m *datadogMonitor) CreateWrapHandleFunc(path string) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tx, ctx := NewDatadogTransaction(r.Context(), path)
			defer tx.End()

			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// Transaction

type datadogTransaction struct {
	finishOptions []ddtrace.FinishOption
	span          ddtrace.Span
}

func NewDatadogTransaction(ctx context.Context, spanName string) (Transaction, context.Context) {
	s, ctx := tracer.StartSpanFromContext(ctx, spanName)
	tx := &datadogTransaction{span: s}
	return tx, ContextWithTransaction(ctx, tx)
}

func (t *datadogTransaction) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return NewDatadogTransaction(ctx, name)
}

func (t *datadogTransaction) End() {
	t.span.Finish(t.finishOptions...)
}

func (t *datadogTransaction) AddAttribute(key string, value interface{}) {
	t.span.SetTag(key, value)
}

func (t *datadogTransaction) AddAttributes(attrs map[string]interface{}) {
	for k, v := range attrs {
		t.AddAttribute(k, v)
	}
}

func (t *datadogTransaction) NoticeError(err error) {
	t.finishOptions = append(t.finishOptions, tracer.WithError(err))
}
