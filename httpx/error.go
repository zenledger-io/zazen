package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zenledger-io/go-utils/logger"
	"go.uber.org/zap"
)

// Error is a JSON-serializable HTTP error.
type Error struct {
	// Err is an optional error that will be logged and written
	// to the wire. If a Message is also supplied, the Err and Message
	// will be formatted as fmt.Sprintf("%s: %s", Message, Err) and the
	// resulting value will be logged and written to the wire.
	Err error

	// Message is an optional message that will be logged and written to
	// the wire. If an Error is also supplied, the Error and Message
	// will be formatted as fmt.Sprintf("%s: %s", Message, Err) and the
	// resulting value will be logged and written to the wire.
	Message string

	// StatusCode is the status code for the error. If this value is not set
	// an HTTP 500 will be written to the wire.
	StatusCode int
}

// errorResponse is an error response.
type errorResponse struct {
	Message string `json:"message"`
}

// Write writes the error to the http.Writer.
func (e *Error) Write(ctx context.Context, w http.ResponseWriter) {
	l := logger.FromContext(ctx)

	w.Header().Set("Content-Type", jsonContentType)

	code := e.StatusCode
	if code <= 0 {
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)

	r := errorResponse{
		Message: http.StatusText(code),
	}

	if e.Err != nil && e.Message != "" {
		r.Message = fmt.Sprintf("%s: %s", e.Err, e.Message)
	} else if e.Err != nil {
		r.Message = e.Err.Error()
	} else if e.Message != "" {
		r.Message = e.Message
	}

	if err := json.NewEncoder(w).Encode(&r); err != nil {
		l.Error("json encode error response", zap.Error(err))
	}
}
