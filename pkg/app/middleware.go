package app

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	// body []byte
}

func middleware(router http.Handler) http.Handler {
	return loggingMiddleware(router)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)
		now := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(now)
		log.Printf("Received Request: %s %s %dms %d", r.Method, r.URL.Path, duration.Milliseconds(), lrw.statusCode)
	})
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
