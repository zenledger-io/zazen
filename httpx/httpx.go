package httpx

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	jsonContentType = "application/json; charset=utf-8"
)

// Writeable is an interface for anything that can
// be written to an http.ResponseWriter.
type Writeable interface {
	GetStatusCode() int
	GetPayload() any
}

// WriteJSON serializes a Writeable to an http.ResponseWriter.
func WriteJSON(ctx context.Context, w http.ResponseWriter, ww Writeable) {
	w.Header().Set("Content-Type", jsonContentType)

	w.WriteHeader(ww.GetStatusCode())

	if err := json.NewEncoder(w).Encode(ww.GetPayload()); err != nil {
		// log error
	}
}
