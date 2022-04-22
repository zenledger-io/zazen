package middlewarex

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zenledger-io/zazen/observe"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				obs := observe.FromContext(r.Context())
				obs.Panic(r.Context(), fmt.Sprintf("%+v", rvr), debug.Stack())

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
