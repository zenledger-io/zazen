package middlewarex

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// ZapLogFormatter is a go-chi/chi/middleware.LogFormatter that
// writes to a *zap.Logger.
type ZapLogFormatter struct {
	Logger *zap.Logger
}

func (l *ZapLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := l.Logger.With(
		zap.Time("ts_utc", time.Now().UTC()),
		zap.String("req_id", middleware.GetReqID(r.Context())),
		zap.String("scheme", r.URL.Scheme),
		zap.String("proto", r.Proto),
		zap.String("method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("url", r.URL.String()),
	)

	logger.Info("begin request")

	return &zapLoggerEntry{
		Logger: logger,
	}
}

type zapLoggerEntry struct {
	Logger *zap.Logger
}

func (l *zapLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.With(
		zap.Int("resp_status", status),
		zap.Int("resp_bytes", bytes),
		zap.Int("resp_elapsed_ms", int(elapsed.Nanoseconds()/1000000)))

	l.Logger.Info("end request")
}

func (l *zapLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)))
}
