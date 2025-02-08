package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// loggingMiddleware logs incoming HTTP requests
func LoggingMiddleware(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Call the next handler in the chain
		duration := time.Since(start)
		if r.URL.Path == "/ready" || r.URL.Path == "/live" {
			return
		}
		logger.Desugar().Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", duration),
		)
	})
}
