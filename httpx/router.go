package httpx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"

	"github.com/zenledger-io/go-utils/httpx/middlewarex"
	"github.com/zenledger-io/go-utils/observe"
	"github.com/zenledger-io/go-utils/version"
)

var (
	defaultRequestTimeout = 60 * time.Second
)

// APIRouterConfig is configuration for an APIRouter.
type APIRouterConfig struct {
	// APIRouters is a map of API versions to subroutes. All routes
	// created on a router with this config will be nested under
	// "/v1" when APIVersion is 1.
	APIRouters map[int]chi.Router

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

	r.Use(middlewarex.Recoverer)

	// Request timeouts
	timeout := cfg.RequestTimeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}
	r.Use(middleware.Timeout(timeout))

	r.Use(chitrace.Middleware(chitrace.WithServiceName(cfg.ServiceName)))

	// Default routes
	r.Get("/status", status)

	// API routers
	for v, rr := range cfg.APIRouters {
		ver := fmt.Sprintf("/v%d", v)
		r.Mount(ver, rr)
	}

	return r
}

// status is a default status handler.
func status(w http.ResponseWriter, r *http.Request) {
	obs := observe.FromContext(r.Context())

	obs.Debug(r.Context(), "status")

	type response struct {
		Version string `json:"version"`
		Hash    string `json:"hash"`
	}

	NewResponseOK(response{
		Version: version.Version,
		Hash:    version.Hash,
	}).Write(r.Context(), w)
}
