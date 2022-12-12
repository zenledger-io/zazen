package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
)

const (
	DefaultFlags = log.LUTC | log.Ldate | log.Ltime | log.Lmsgprefix

	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var (
	reporter     = NewNoneReporter()
	zapLogger, _ = func() (*zap.Logger, error) {
		cfg := zap.NewProductionConfig()
		return cfg.Build()
	}()
)

func init() {
	log.SetFlags(DefaultFlags)
}

func Start(ctx context.Context) {
	var hbReporter honeybadgerReporter
	if err := hbReporter.Start(ctx); err != nil {
		ContextLogger(ctx).Printf("%v", err)
	} else {
		reporter = hbReporter
	}
}

func Printf(format string, args ...interface{}) {
	zapLogger.Info(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	if err := reporter.Errorf(format, args...); err != nil {
		zapLogger.Error(fmt.Sprintf("reporting error: %v", err))
	}

	zapLogger.Error(fmt.Sprintf(format, args...))
}
