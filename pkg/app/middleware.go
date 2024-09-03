package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	// body []byte
}

type contextKey string

const userInfoKey contextKey = "userInfo"

func middleware(router http.Handler) http.Handler {
	return loggingMiddleware(JWTAuthMiddleware(router))
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

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/users/login" {
			authorizationHeader := r.Header.Get("Authorization")

			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

			if !strings.Contains(authorizationHeader, "Bearer") {
				http.Error(w, "Invalid Token", http.StatusBadRequest)
				return
			}

			claims, err := parseAuthToken(tokenString)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(context.Background(), userInfoKey, claims)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
