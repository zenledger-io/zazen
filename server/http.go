package server

import (
	"github.com/zenledger-io/go-utils/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	DefaultMiddleware = []middleware.Func{
		middleware.Logger(),
		middleware.UUID(),
	}
)

type HTTPEndpoint struct {
	Path       string
	Methods    []string
	Handler    http.HandlerFunc
	Middleware []middleware.Func
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
		r.Handle(e.Path, middleware.Wrap(e.Handler, e.Middleware...)).Methods(e.Methods...)
	}

	return middleware.Wrap(r.ServeHTTP, defaultMiddleware...)
}
