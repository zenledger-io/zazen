package server

import (
	"github.com/zenledger-io/zazen/middleware"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
)

var (
	DefaultMiddleware = []middleware.Func{
		middleware.Logger(),
		middleware.UUID(),
	}
	MetricsMiddleware func(path string) middleware.Func
)

type HTTPEndpoint struct {
	Path       string
	Methods    []string
	Handler    http.HandlerFunc
	Middleware []middleware.Func
}

type HTTPEndpointGroup struct {
	GroupPath  string
	Endpoints  []HTTPEndpoint
	Middleware []middleware.Func
}

func (g HTTPEndpointGroup) MakeEndpoints() []HTTPEndpoint {
	for i, e := range g.Endpoints {
		e.Path = path.Join(g.GroupPath, e.Path)
		e.Middleware = append(e.Middleware, g.Middleware...)
		g.Endpoints[i] = e
	}

	return g.Endpoints
}

func NewHTTP(endpoints []HTTPEndpoint, addr string, rTimeout, wTimeout time.Duration) *http.Server {
	srv := http.Server{
		Addr:         addr,
		Handler:      Router(endpoints, DefaultMiddleware...),
		ReadTimeout:  rTimeout,
		WriteTimeout: wTimeout,
	}
	return &srv
}

func Router(endpoints []HTTPEndpoint, defaultMiddleware ...middleware.Func) http.HandlerFunc {
	r := mux.NewRouter()
	for _, e := range endpoints {
		e.Middleware = append(e.Middleware, defaultMiddleware...)
		if MetricsMiddleware != nil {
			e.Middleware = append(e.Middleware, MetricsMiddleware(e.Path))
		}
		r.Handle(e.Path, middleware.Wrap(e.Handler, e.Middleware...)).Methods(e.Methods...)
	}

	return r.ServeHTTP
}
