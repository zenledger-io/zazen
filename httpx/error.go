package httpx

import (
	"fmt"
	"net/http"
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

// GetStatusCode returns the status code for the error.
func (e Error) GetStatusCode() int {
	if e.StatusCode <= 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

// GetPayload returns the payload for the error.
func (e Error) GetPayload() any {
	r := errorResponse{
		Message: http.StatusText(e.GetStatusCode()),
	}

	if e.Err != nil && e.Message != "" {
		r.Message = fmt.Sprintf("%s: %s", e.Err, e.Message)
	} else if e.Err != nil {
		r.Message = e.Err.Error()
	} else if e.Message != "" {
		r.Message = e.Message
	}

	return &r
}

type errorResponse struct {
	Message string `json:"message"`
}
