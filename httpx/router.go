package httpx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/zenledger-io/go-utils/httpx/middlewarex"
	"github.com/zenledger-io/go-utils/version"
)

var (
	defaultRequestTimeout = 60 * time.Second
)

// APIRouterConfig is configuration for an APIRouter.
type APIRouterConfig struct {
	// APIVersion is the URL-based API version. All routes created
	// on a router with this config will be nested under
	// "/v1" when APIVersion is 1.
	APIVersion int

	// APIRouter is the router to be mounted under the API verison.
	APIRouter chi.Router

	// Logger is the base logger to be used for all requests.
	Logger *zap.Logger

	// RequestTimeout is the maximum amount of time allowed
	// for a request on the router before processing is stopped.
	RequestTimeout time.Duration

	// ServiceName is the name of the service running this router.
	ServiceName string
}

// NewAPIRouter creates a chi Router that is configured with a set
// of default middleware and routes for an API. This router reads
// and writes JSON.
func NewAPIRouter(cfg APIRouterConfig) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Use(middleware.RequestLogger(&middlewarex.ZapLogFormatter{
		Logger: cfg.Logger,
	}))

	r.Use(middleware.Recoverer)

	// Request timeouts
	timeout := cfg.RequestTimeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}
	r.Use(middleware.Timeout(timeout))

	// Prometheus
	r.Use(middlewarex.NewMetrics())

	// Default routes
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/status", status)

	// API
	ver := fmt.Sprintf("/v%d", cfg.APIVersion)
	r.Mount(ver, cfg.APIRouter)

	return r
}

// status is a default status handler.
func status(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Version string `json:"version"`
		Hash    string `json:"hash"`
	}

	NewResponseOK(response{
		Version: version.Version,
		Hash:    version.Hash,
	}).Write(r.Context(), w)
}
