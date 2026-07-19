package middleware

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &StatusRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		next.ServeHTTP(rec, r) 

		log.Printf(
			"%s | %s %s | %d | %s",
			start.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			rec.status,
			time.Since(start),
		)
	})
}
