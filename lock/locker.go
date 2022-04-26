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
