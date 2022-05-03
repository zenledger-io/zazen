package httpx

import (
	"net/http"
)

// Response is a JSON-serializable HTTP response.
type Response struct {
	// Payload is a JSON-serializable value that
	// will be serialized to the wire.
	Payload interface{}

	// StatusCode is the status code for the response.
	StatusCode int
}

// GetStatusCode returns the status code for the response.
func (r Response) GetStatusCode() int {
	if r.StatusCode <= 0 {
		return http.StatusOK
	}
	return r.StatusCode
}

// GetPayload returns the payload for the response.
func (r Response) GetPayload() any {
	return r.Payload
}
