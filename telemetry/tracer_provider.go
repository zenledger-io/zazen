package telemetry

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TracerProviderConfig is configuration for a trace provider.
type TracerProviderConfig struct {
	// TargetAddr is the address of the OTLP collector to which traces are sent. If this
	// value is not set, NewTraceProvider will write to Writer.
	TargetAddr string

	// If this value is set and TargetAddr is left blank then traces wiil be written to it.
	// If this value is also empty, traces will be written to iotil.Discard.
	TargetWriter io.Writer

	// TransportCredentials are credentials to use when dialing TargetAddr.
	TransportCredentials credentials.TransportCredentials

	// Sampler is use to inform the tracer provider as to how
	// sampling should be handled in the spans it creates.
	Sampler sdktrace.Sampler
}

// NewTracerProvider creates a new OTEL trace provider.
func NewTracerProvider(ctx context.Context, cfg Config, tpCfg TracerProviderConfig) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String(cfg.Name),
		semconv.ServiceVersionKey.String(cfg.BuildVersion),
		attribute.Key("service.hash").String(cfg.BuildVersion),
	))
	if err != nil {
		return nil, fmt.Errorf("new resource: %w", err)
	}

	var exporter sdktrace.SpanExporter
	if tpCfg.TargetAddr == "" {
		exporter, err = newExporterStdout(tpCfg)
	} else {
		exporter, err = newExporterGRPC(ctx, tpCfg)
	}
	if err != nil {
		return nil, fmt.Errorf("new exporter: %w", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(tpCfg.Sampler)), nil
}

func newExporterGRPC(ctx context.Context, cfg TracerProviderConfig) (sdktrace.SpanExporter, error) {
	conn, err := grpc.DialContext(ctx, cfg.TargetAddr,
		grpc.WithTransportCredentials(cfg.TransportCredentials),
		grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("grpc dial context: %w", err)
	}

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("new grpc trace exporter: %w", err)
	}

	return exp, nil
}

func newExporterStdout(cfg TracerProviderConfig) (sdktrace.SpanExporter, error) {
	w := ioutil.Discard
	if cfg.TargetWriter != nil {
		w = cfg.TargetWriter
	}
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint())
}
