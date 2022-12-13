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

func Sync() {
	_ = zapLogger.Sync()
}

func Debugf(format string, args ...any) {
	zapLogger.Debug(fmt.Sprintf(format, args...))
}

func Printf(format string, args ...any) {
	zapLogger.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	zapLogger.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	if err := reporter.Errorf(format, args...); err != nil {
		zapLogger.Error(fmt.Sprintf("reporting error: %v", err))
	}

	zapLogger.Error(fmt.Sprintf(format, args...))
}

func DebugT(msg string, tags ...Tag) {
	zapLogger.Debug(msg, tagsToFields(tags)...)
}

func PrintT(msg string, tags ...Tag) {
	zapLogger.Info(msg, tagsToFields(tags)...)
}

func WarnT(msg string, tags ...Tag) {
	zapLogger.Warn(msg, tagsToFields(tags)...)
}

func ErrorT(msg string, tags ...Tag) {
	if err := reporter.Errorf(msg); err != nil {
		zapLogger.Error(fmt.Sprintf("reporting error: %v", err))
	}

	zapLogger.Error(msg, tagsToFields(tags)...)
}

func tagsToFields(tags []Tag) []zap.Field {
	fields := make([]zap.Field, len(tags))
	for i, t := range tags {
		fields[i] = t.field()
	}
	return fields
}
