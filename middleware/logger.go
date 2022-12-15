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

			log.ContextLogger(r.Context()).PrintT("finished http request",
				log.NewTag("duration", fmt.Sprintf("%v", time.Since(t))),
				log.NewTag("http.status_code", wWrapper.StatusCode),
				log.NewTag("http.url", r.URL.Path),
				log.NewTag("http.method", r.Method),
				log.NewTag("http.request.bytes", bWrapper.ByteLength),
				log.NewTag("http.response.bytes", wWrapper.ByteLength))
		}
	}
}
