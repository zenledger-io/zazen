package service

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/trace"

	"github.com/zenledger-io/zazen/httpx"
	"github.com/zenledger-io/zazen/telemetry"
)

var (
	defaultRequestTimeout = 60 * time.Second
)

func newRoot(tp trace.TracerProvider, cfg Config) http.Handler {
	root := chi.NewRouter()

	timeout := cfg.RequestTimeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}

	root.Use(middleware.Timeout(timeout))
	root.Use(middleware.RealIP)
	root.Use(middleware.RequestID)

	root.Use(otelchi.Middleware(cfg.TelemetryConfig.Name,
		otelchi.WithTracerProvider(tp)))
	root.Use(telemetry.RecoverMiddleware)

	root.Get("/status", newStatusHandler(cfg.TelemetryConfig))
	for p, h := range cfg.Mounts {
		root.Mount(p, h)
	}

	return root
}

func newStatusHandler(telCfg telemetry.Config) http.HandlerFunc {
	type response struct {
		Version string `json:"version"`
		Hash    string `json:"hash"`
	}

	payload := response{
		Version: telCfg.BuildVersion,
		Hash:    telCfg.BuildHash,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(r.Context(), w, httpx.Response{
			Payload: payload,
		})
	}
}
