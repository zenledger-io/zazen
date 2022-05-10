package lock

import (
	"context"
)

type Locker interface {
	Lock(context.Context) (Unlocker, error)
}

type Unlocker interface {
	Unlock() error
}

type RefreshLocker interface {
	Locker
	LockWithRefresh(context.Context) (RefreshUnlocker, error)
	Indefinite(ctx context.Context, f func(context.Context)) error
}

type RefreshUnlocker interface {
	Unlocker
	Refresh(context.Context) error
}
