package retry

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zenledger-io/zazen/errors"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tcs := map[string]struct {
		Config          Config
		ContextFunc     func(context.Context) context.Context
		Func            func(context.Context) error
		ErrExpected     bool
		RetriesExpected int
	}{
		"with max retries": {
			Config: Config{
				ShouldRetry: MaxRetries(3),
			},
			ContextFunc: func(ctx context.Context) context.Context {
				return ctx
			},
			Func: func(ctx context.Context) error {
				return errors.Errorf("some error")
			},
			ErrExpected:     true,
			RetriesExpected: 3,
		},
		"with max retries and error check": {
			Config: Config{
				ShouldRetry: WrapShouldRetryFuncs(MaxRetries(3), func(ctx context.Context, i int, err error) bool {
					if err.Error() == "stop" {
						return false
					}

					return true
				}),
			},
			ContextFunc: func(ctx context.Context) context.Context {
				return ctx
			},
			Func: func() func(context.Context) error {
				var i int
				return func(ctx context.Context) error {
					i += 1
					if i > 2 {
						return errors.Errorf("stop")
					}

					return errors.Errorf("some error")
				}
			}(),
			ErrExpected:     true,
			RetriesExpected: 2,
		},
		"with max retries and reset count after": {
			Config: Config{
				ShouldRetry: WrapShouldRetryFuncs(MaxRetries(3)),
				ResetCount:  ResetCountAfter(50 * time.Millisecond),
			},
			ContextFunc: func(ctx context.Context) context.Context {
				return ctx
			},
			Func: func() func(context.Context) error {
				var i int
				return func(ctx context.Context) error {
					i += 1
					switch i {
					case 2:
						time.Sleep(50 * time.Millisecond)
					case 3:
						return nil
					}

					return errors.Errorf("some error")
				}
			}(),
			ErrExpected:     false,
			RetriesExpected: 1,
		},
		"context error returns immediately": {
			Config: Config{
				ShouldRetry: MaxRetries(3),
			},
			ContextFunc: func(ctx context.Context) context.Context {
				ctx, cancel := context.WithCancel(ctx)
				cancel()
				return ctx
			},
			Func: func(ctx context.Context) error {
				return errors.Errorf("some error")
			},
			ErrExpected:     true,
			RetriesExpected: 0,
		},
		"indefinite retry": {
			Config: Config{},
			ContextFunc: func(ctx context.Context) context.Context {
				return ctx
			},
			Func: func() func(context.Context) error {
				var i int
				return func(ctx context.Context) error {
					i += 1
					if i > 5 {
						return nil
					}

					return errors.Errorf("some error")
				}
			}(),
			ErrExpected:     false,
			RetriesExpected: 5,
		},
		"indefinite retry context error returns immediately": {
			Config: Config{},
			ContextFunc: func(ctx context.Context) context.Context {
				ctx, cancel := context.WithCancel(ctx)
				cancel()
				return ctx
			},
			Func: func(ctx context.Context) error {
				return errors.Errorf("some error")
			},
			ErrExpected:     true,
			RetriesExpected: 0,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			i, err := tc.Config.Do(tc.ContextFunc(ctx), tc.Func)
			if tc.ErrExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.RetriesExpected, i)
		})
	}
}
