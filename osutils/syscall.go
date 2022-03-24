package osutils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	DefaultSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
	}
)

func WatchSignals(ctx context.Context, cancel context.CancelFunc, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultSignals
	}

	stop := make(chan os.Signal, len(signals))
	for _, s := range signals {
		signal.Notify(stop, s)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			cancel()
		}
	}
}
