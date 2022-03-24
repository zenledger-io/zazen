package middleware

import (
	"github.com/google/uuid"
	"github.com/zenledger-io/go-utils/log"
	"net/http"
)

var (
	UUIDHeader string
)

func UUID() Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var rid string
			if UUIDHeader != "" {
				rid = r.Header.Get(UUIDHeader)
			}

			if rid == "" {
				rid = uuid.New().String()
			}

			ctx := log.ContextWithUUID(r.Context(), rid)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
