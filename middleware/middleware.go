package middleware

import "net/http"

// Wrap calls generates a new handler out of the middleware passed in.
// Middleware gets executed last to first.
func Wrap(h http.HandlerFunc, mw ...Func) http.HandlerFunc {
	for _, m := range mw {
		h = m(h)
	}

	return h
}
