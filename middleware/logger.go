package middleware

import (
	"github.com/zenledger-io/zazen/httputils"
	"github.com/zenledger-io/zazen/ioutils"
	"github.com/zenledger-io/zazen/log"
	"net/http"
	"time"
)

func Logger() Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			wWrapper := httputils.NewMeasuredResponseWriter(w)
			bWrapper := ioutils.NewMeasuredReadCloser(r.Body)

			r.Body = bWrapper
			next.ServeHTTP(wWrapper, r)

			log.ContextLogger(r.Context()).Printf("%v %v %v | %v | %v | sent %vB | received %vB",
				wWrapper.StatusCode,
				r.Method,
				r.URL,
				r.RemoteAddr,
				time.Since(t),
				wWrapper.ByteLength,
				bWrapper.ByteLength)
		}
	}
}
