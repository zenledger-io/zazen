package metrics

import (
	"context"
	"fmt"
	"github.com/zenledger-io/zazen/log"
	"net/http"
)

// Monitor

func NewStdOutMonitor() Monitor {
	return stdOutMonitor{}
}

type stdOutMonitor struct{}

func (m stdOutMonitor) Start(context.Context) error {
	return nil
}

func (m stdOutMonitor) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return NewStdOutTransaction(ctx, name)
}

func (m stdOutMonitor) CreateWrapHandleFunc(name string) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tx, ctx := NewStdOutTransaction(r.Context(), name)
			defer tx.End()

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		}
	}
}

// Transaction

type stdOutTransaction struct {
	depth  int
	name   string
	logger log.Logger
	attrs  map[string]any
}

func NewStdOutTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	logger := log.ContextLogger(ctx)
	var depth int
	if tx, ok := ContextTransaction(ctx); ok {
		if stx, ok := tx.(*stdOutTransaction); ok {
			depth = stx.depth + 1
		}
	}
	logger.PrintT(fmt.Sprintf("starting %v", name), log.NewTag("depth", depth))
	tx := &stdOutTransaction{
		logger: logger,
		depth:  depth,
		name:   name,
		attrs:  make(map[string]any),
	}
	ctx = ContextWithTransaction(ctx, tx)
	return tx, ctx
}

func (t *stdOutTransaction) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return NewStdOutTransaction(ctx, name)
}

func (t *stdOutTransaction) End() {
	tags := make([]log.Tag, 0, len(t.attrs)+1)
	tags = append(tags, log.NewTag("depth", t.depth))
	for k, v := range t.attrs {
		tags = append(tags, log.NewTag(k, v))
	}
	t.logger.PrintT(fmt.Sprintf("finished %v", t.name), tags...)
}

func (t *stdOutTransaction) AddAttribute(key string, value any) {
	t.attrs[key] = value
}

func (t *stdOutTransaction) AddAttributes(attrs map[string]any) {
	for k, v := range attrs {
		t.attrs[k] = v
	}
}

func (t *stdOutTransaction) NoticeError(err error) {
	t.attrs["error"] = err
}
