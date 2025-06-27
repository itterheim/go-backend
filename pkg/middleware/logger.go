package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := time.Now().Format(time.RFC3339)

		start := time.Now()
		log.Printf("Started %s %s %s", requestID, r.Method, r.URL.Path)

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("Completed %s %s in %v with status %d", requestID, r.URL.Path, time.Since(start), rw.statusCode)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}