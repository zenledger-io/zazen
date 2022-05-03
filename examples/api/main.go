package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/zenledger-io/zazen/httpx"
	"github.com/zenledger-io/zazen/service"
	"github.com/zenledger-io/zazen/telemetry"
)

const (
	serviceName    = "the-birds"
	serviceVersion = "0.0.1"
	addr           = ":8111"
)

type bird struct {
	Name string `json:"name"`
}

var birds = []bird{
	{Name: "eagle"},
	{Name: "hawk"},
	{Name: "pidgeon"},
	{Name: "stork"},
}

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		failf("failed to create service: %v\n", err)
	}

	defer l.Sync()

	ctx := telemetry.ContextWithLog(context.Background(), telemetry.NewZapLog(l))

	svc, err := service.New(ctx, service.Config{
		Name:         serviceName,
		BuildVersion: serviceVersion,
		BuildHash:    "abc111",
		TracerProviderConfig: telemetry.TracerProviderConfig{
			ServiceName:    serviceName,
			ServiceVersion: serviceVersion,
			// To see traces, uncomment the following line
			// TargetWriter: os.Stderr,
		},
		Mounts: map[string]http.Handler{
			"/v1": v1Handler(),
		},
	})
	if err != nil {
		failf("failed to create service: %v\n", err)
	}

	defer func() {
		if err := svc.Shutdown(context.Background()); err != nil {
			failf("failed to shutdown service: %v", err)
		}
	}()

	server := &http.Server{
		Addr:    addr,
		Handler: svc.Handler(),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	l.Info("listening", zap.String("addr", addr))
	server.ListenAndServe()
}

func v1Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/birds", func(w http.ResponseWriter, r *http.Request) {
		span := telemetry.SpanFromContext(r.Context())
		span.Debug("inside birds handler")

		time.Sleep(100 * time.Millisecond)

		httpx.WriteJSON(r.Context(), w, httpx.Response{
			Payload: birds,
		})
	})

	r.Get("/boom", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	return r
}

func failf(msgFmt string, err error) {
	fmt.Printf(msgFmt, err)
	os.Exit(1)
}
