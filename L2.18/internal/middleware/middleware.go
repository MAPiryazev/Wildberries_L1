package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware для того чтобы логировать запросы, идущие через handler
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(lrw, r)

		fmt.Printf("[%s] %s %s %d %s\n",
			start.Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			lrw.status,
			time.Since(start),
		)
	})
}
