package httputils

import (
	"context"
	"encoding/json"
	"github.com/zenledger-io/zazen/log"
	"net/http"
)

var (
	JSONContentType = "application/json; charset=utf-8"
)

func SendJSON[T any](w http.ResponseWriter, t T) error {
	w.Header().Set("Content-Type", JSONContentType)
	return json.NewEncoder(w).Encode(t)
}

func SendError[T any](ctx context.Context, statusCode int, w http.ResponseWriter, t T) {
	w.WriteHeader(statusCode)
	if err := SendJSON(w, t); err != nil {
		log.ContextLogger(ctx).Errorf("unable to send json: %v", err)
	}
}
