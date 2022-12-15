package middleware

import (
	"fmt"
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

			dur := time.Since(t)
			log.ContextLogger(r.Context()).PrintT(
				fmt.Sprintf("finished http request %v %v", r.Method, r.URL.Path),
				log.NewTag("duration", dur.Nanoseconds()),
				log.NewTag("duration_formatted", fmt.Sprintf("%v", dur)),
				log.NewTag("http.status_code", wWrapper.StatusCode),
				log.NewTag("http.path", r.URL.Path),
				log.NewTag("http.method", r.Method),
				log.NewTag("http.request.bytes", bWrapper.ByteLength),
				log.NewTag("http.response.bytes", wWrapper.ByteLength))
		}
	}
}
