package metrics

import (
	"context"
	"fmt"
	"net/http"
)

type Monitor interface {
	Start(context.Context) error
	StartTransaction(ctx context.Context, name string) (Transaction, context.Context)
	CreateWrapHandleFunc(path string) func(h http.HandlerFunc) http.HandlerFunc
}

func Start(ctx context.Context) (context.Context, error) {
	var monitor Monitor
	switch MonitorType {
	case "datadog":
		monitor = NewDatadogMonitor()
	default:
		monitor = NewNullMonitor()
	}

	if err := monitor.Start(ctx); err != nil {
		wrappedErr := fmt.Errorf("failed to start %v monitor: %w", MonitorType, err)
		return ContextWithMonitor(ctx, NewNullMonitor()), wrappedErr
	}

	return ContextWithMonitor(ctx, monitor), nil
}
