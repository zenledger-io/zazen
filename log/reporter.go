package log

import "context"

type Reporter interface {
	Start(ctx context.Context) error
	Errorf(string, ...interface{}) error
	Monitor()
}
