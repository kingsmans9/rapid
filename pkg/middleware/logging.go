package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/spectrocloud/rapid-agent/pkg/util/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logWriter := NewLoggingResponseWriter(w)
		next.ServeHTTP(logWriter, r)

		if os.Getenv("DEBUG") != "true" && logWriter.StatusCode < http.StatusBadRequest {
			return
		}

		logger.Infof(
			"method=%s status=%d duration=%s request=%s",
			r.Method,
			logWriter.StatusCode,
			time.Since(startTime).String(),
			r.RequestURI,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}
