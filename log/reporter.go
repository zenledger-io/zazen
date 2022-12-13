package log

import "context"

type Reporter interface {
	Start(ctx context.Context) error
	Errorf(string, ...any) error
	Monitor()
}
