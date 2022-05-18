package telemetry

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLog returns a Log that carries a zap logger..
func NewZapLog(logger *zap.Logger) Log {
	return &otelZapLog{
		logger: otelzap.New(logger, otelzap.WithMinLevel(zapcore.DebugLevel)),
	}
}

// otelZapLog wraps an OTEL-aware zap logger.
// See https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap
type otelZapLog struct {
	logger *otelzap.Logger
}

func (z *otelZapLog) Debug(ctx context.Context, message string, fields ...Field) {
	z.logger.Ctx(ctx).Debug(message, zapFields(fields, 0)...)
}

func (z *otelZapLog) Info(ctx context.Context, message string, fields ...Field) {
	z.logger.Ctx(ctx).Info(message, zapFields(fields, 0)...)
}

func (z *otelZapLog) Error(ctx context.Context, message string, err error, fields ...Field) {
	zfs := append(zapFields(fields, 1), zap.Error(err))
	z.logger.Ctx(ctx).Error(message, zfs...)
}

func zapFields(fields []Field, padding int) []zap.Field {
	zf := make([]zap.Field, 0, len(fields)+padding)
	for _, f := range fields {
		zf = append(zf, zap.Any(f.Name, f.Value))
	}
	return zf
}
