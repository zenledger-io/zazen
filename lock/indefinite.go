package lock

import (
	"context"
	"github.com/zenledger-io/zazen/timeutils"
	"time"
)

func Indefinite(ctx context.Context, locker RefreshLocker, interval time.Duration, f func(context.Context)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	l, err := locker.LockWithRefresh(ctx)
	if err != nil {
		return err
	}
	defer l.Unlock()

	go func() {
		defer cancel()

		f(ctx)
	}()

	t := time.NewTicker(interval)
	defer timeutils.StopTicker(t)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if err := l.Refresh(ctx); err != nil {
				return err
			}
		}
	}
}
