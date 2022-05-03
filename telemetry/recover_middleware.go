package telemetry

import (
	"fmt"
	"net/http"
	"runtime"
)

// RecoverMiddleware catches panics and logs them to a span.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			span := SpanFromContext(r.Context())
			if !span.IsRecording() {
				return
			}
			if err := recover(); err != nil && err != http.ErrAbortHandler {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				span.Error("panic", fmt.Errorf("%v\n%s", err, buf))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
