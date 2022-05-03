package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext returns new Span based on values
// carrioed in ctx.
func SpanFromContext(ctx context.Context) Span {
	return &span{
		Span: trace.SpanFromContext(ctx),
		log:  LogFromContext(ctx),
		ctx:  ctx,
	}
}

// Span is the most granluar component of a trace.
type Span interface {
	trace.Span
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Error(message string, err error, fields ...Field)
}

type span struct {
	trace.Span
	log Log
	ctx context.Context
}

func (s *span) Debug(message string, fields ...Field) {
	s.log.Debug(s.ctx, message, fields...)
}

func (s *span) Info(message string, fields ...Field) {
	s.log.Info(s.ctx, message, fields...)
}

func (s *span) Error(message string, err error, fields ...Field) {
	s.log.Error(s.ctx, message, err, fields...)
	s.Span.RecordError(err)
}
