package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zenledger-io/zazen/telemetry"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config is configuration for a Service.
type Config struct {
	// TelemetryConfig is the configuration that will be used
	// to instrument this service.
	TelemetryConfig telemetry.Config

	// RequestTimeout is the maximum amount of time allowed
	// for a request on the router before processing is stopped.
	RequestTimeout time.Duration

	// Mounts are a map of paths to http.Handler that will be
	// loaded under the root of the service.
	Mounts map[string]http.Handler

	// TracerProviderConfig is configuration for the trace provider
	// that the service will create.
	TracerProviderConfig telemetry.TracerProviderConfig
}

// New creates a new Service.
func New(ctx context.Context, cfg Config) (*Service, error) {
	tp, err := telemetry.NewTracerProvider(ctx, cfg.TelemetryConfig, cfg.TracerProviderConfig)
	if err != nil {
		return nil, fmt.Errorf("new trace provider: %w", err)
	}

	root := newRoot(tp, cfg)

	return &Service{
		root: root,
		tp:   tp,
	}, nil
}

// Service is a component used to serve versioned,
// HTTP services that speak JSON.
type Service struct {
	root http.Handler
	tp   *sdktrace.TracerProvider
}

// Handler returns the root handler for the service.
func (s *Service) Handler() http.Handler {
	return s.root
}

// TracerProvider returns the OTEL trace provider being
// used by the service.
func (s *Service) TracerProvider() *sdktrace.TracerProvider {
	return s.tp
}

// Shutdown shuts down a service.
func (s *Service) Shutdown(ctx context.Context) error {
	if err := s.tp.Shutdown(ctx); err != nil {
		return fmt.Errorf("tracer provider shutdown: %w", err)
	}

	return nil
}
