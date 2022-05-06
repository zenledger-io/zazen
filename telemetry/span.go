package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext returns the existing Span based on values
// carrioed in ctx.
func SpanFromContext(ctx context.Context) Span {
	return &span{
		Span: trace.SpanFromContext(ctx),
		log:  LogFromContext(ctx),
		ctx:  ctx,
	}
}

// Start creates a span and a context.Context containing the newly-created span
// using the supplied tracer and options.
func Start(ctx context.Context, tracer trace.Tracer, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	var span trace.Span
	ctx, span = tracer.Start(name, opts...)

	return ctx, &span{
		Span: span,
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
	s.Span.SetStatus(codes.Error, err.Error())
}
