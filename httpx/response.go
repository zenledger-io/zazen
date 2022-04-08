package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zenledger-io/go-utils/logger"
	"go.uber.org/zap"
)

const (
	jsonContentType = "application/json; charset=utf-8"
)

// Response is a JSON-serializable HTTP response.
type Response struct {
	// Payload is a JSON-serializable value that
	// will be serialized to the wire.
	Payload interface{}

	// StatusCode is the status code for the response.
	StatusCode int
}

// NewResponse creates a new Response for a payload
// with HTTP 200 as the status code.
func NewResponse(payload interface{}) *Response {
	return &Response{
		Payload:    payload,
		StatusCode: http.StatusOK,
	}
}

// Write writes a response to an http.ResponseWriter.
// Any errors encountered will be automatically handled.
func (r *Response) Write(ctx context.Context, w http.ResponseWriter) {
	l := logger.FromContext(ctx)

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(r.StatusCode)

	if err := json.NewEncoder(w).Encode(r.Payload); err != nil {
		l.Error("json encode response", zap.Error(err))
	}
}
