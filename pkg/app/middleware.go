package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adiyakaihsan/go-http-server/pkg/config"
	jwt "github.com/golang-jwt/jwt/v4"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	// body []byte
}

func middleware(router http.Handler) http.Handler {
	return runThroughMiddleware(router)
}

func runThroughMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)

		if r.URL.Path == "/v1/users/login" {
			now := time.Now()
			next.ServeHTTP(lrw, r)
			duration := time.Since(now)
			log.Printf("Received Request: %s %s %dms %d", r.Method, r.URL.Path, duration.Milliseconds(), lrw.statusCode)
		} else {
			authorizationHeader := r.Header.Get("Authorization")

			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	
	
			if !strings.Contains(authorizationHeader, "Bearer"){
				http.Error(w, "Invalid Token", http.StatusBadRequest)
				return
			}
	
			claims, err := parseAuthToken(tokenString)
	
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
	
			ctx := context.WithValue(context.Background(), "userInfo", claims)
	
			r = r.WithContext(ctx)
	
			now := time.Now()
			next.ServeHTTP(lrw, r)
			duration := time.Since(now)
			log.Printf("Received Request: %s %s %dms %d", r.Method, r.URL.Path, duration.Milliseconds(), lrw.statusCode)
		}
	})
}

func parseAuthToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(config.JWT_SIGNATURE_KEY), nil
	})

	if err != nil {
		log.Printf("%s: %v", "Error when parsing token", err)
		return nil,fmt.Errorf("error when parsing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("%s: %v", "Error when validating token", err)
		return nil, fmt.Errorf("error when validating token")
	}

	return claims, nil


}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
