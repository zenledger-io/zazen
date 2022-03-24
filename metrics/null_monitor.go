package metrics

import (
	"context"
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

// Transaction

type nullTransaction struct{}

func (t nullTransaction) StartTransaction(ctx context.Context, name string) (Transaction, context.Context) {
	return nullTransaction{}, ctx
}

func (t nullTransaction) End() {}

func (t nullTransaction) AddAttribute(key string, value interface{}) {}

func (t nullTransaction) AddAttributes(attrs map[string]interface{}) {}

func (t nullTransaction) NoticeError(err error) {}
