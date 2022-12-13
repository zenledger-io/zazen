package metrics

import (
	"context"
	"net/http"
)

// Monitor

type nullMonitor struct{}

func NewNullMonitor() Monitor {
	return nullMonitor{}
}

func (m nullMonitor) Start(context.Context) error {
	return nil
}

func (m nullMonitor) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return nullTransaction{}, ctx
}

func (m nullMonitor) CreateWrapHandleFunc(string) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return h
	}
}

// Transaction

type nullTransaction struct{}

func (t nullTransaction) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return nullTransaction{}, ctx
}

func (t nullTransaction) End() {}

func (t nullTransaction) AddAttribute(key string, value any) {}

func (t nullTransaction) AddAttributes(attrs map[string]any) {}

func (t nullTransaction) NoticeError(err error) {}
