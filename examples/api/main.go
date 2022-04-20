package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zenledger-io/go-utils/httpx"
	"github.com/zenledger-io/go-utils/logger"
	"go.uber.org/zap"
)

const addr = ":8111"

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
		panic(err)
	}

	defer l.Sync()

	l.Info("listening", zap.String("addr", addr))

	http.ListenAndServe(addr, httpx.NewAPIRouter(httpx.APIRouterConfig{
		APIRouters:  newAPIMap(),
		ServiceName: "birds",
		Logger:      l,
	}))
}

func newAPIMap() map[int]chi.Router {
	r := chi.NewRouter()

	r.Get("/birds", func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		l.Info("birds handler")
		time.Sleep(100 * time.Millisecond)
		httpx.NewResponseOK(birds).Write(r.Context(), w)
	})

	r.Get("/boom", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	return map[int]chi.Router{1: r}
}
