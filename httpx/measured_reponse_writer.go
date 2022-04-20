package httpx

import (
	"net/http"
)

type MeasuredResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	ByteLength int
}

func NewMeasuredResponseWriter(w http.ResponseWriter) *MeasuredResponseWriter {
	return &MeasuredResponseWriter{ResponseWriter: w}
}

func (w *MeasuredResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *MeasuredResponseWriter) Write(b []byte) (int, error) {
	i, err := w.ResponseWriter.Write(b)

	// If bytes are sent before writing a header the http package writes status ok header by default.
	if w.StatusCode == 0 && err != http.ErrHijacked { // ErrHijacked bypasses writing the header
		w.StatusCode = http.StatusOK
	}
	w.ByteLength += i

	return i, err
}
