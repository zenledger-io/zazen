package httputils

import (
	"context"
	"encoding/json"
	"github.com/zenledger-io/zazen/log"
	"net/http"
)

const (
	JsonContentType = "application/json; charset=utf-8"
)

func SendJson[T any](w http.ResponseWriter, statusCode int, t T) error {
	w.Header().Set("Content-Type", JsonContentType)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(t)
}

func SendJsonOk[T any](w http.ResponseWriter, t T) error {
	return SendJson(w, http.StatusOK, t)
}

func SendJsonError[T any](ctx context.Context, w http.ResponseWriter, statusCode int, t T) {
	if err := SendJson(w, statusCode, t); err != nil {
		log.ContextLogger(ctx).Errorf("unable to send json: %v", err)
	}
}
