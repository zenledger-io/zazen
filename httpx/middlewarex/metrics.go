package middlewarex

import (
	"net/http"

	metricsprom "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsmware "github.com/slok/go-http-metrics/middleware"
	metricsmwarestd "github.com/slok/go-http-metrics/middleware/std"
)

func NewMetrics() func(next http.Handler) http.Handler {
	m := metricsmware.New(metricsmware.Config{
		Recorder: metricsprom.NewRecorder(metricsprom.Config{}),
	})
	return metricsmwarestd.HandlerProvider("", m)
}
