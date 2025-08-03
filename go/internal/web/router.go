package web

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"eaglebank/internal/reqctx"
)

func NewServer(logger *slog.Logger, validate *validator.Validate, usrSvc UserService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)

	// User routes
	mux.HandleFunc("POST /v1/users", handleCreateUser(validate, usrSvc))

	handler := panicMiddleware(logger)(mux)
	handler = loggingMiddleware(logger)(handler)

	return handler
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.NewString()

			ctx := context.WithValue(r.Context(), reqctx.RequestIDKey, requestID)
			r = r.WithContext(ctx)

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			rw.Header().Set("X-Request-ID", requestID)

			logger.Info("request started",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("referer", r.Referer()),
			)

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			logger.Info("request completed",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status_code", rw.statusCode),
				slog.Int("response_size", rw.size),
				slog.String("duration", duration.String()),
				slog.Float64("duration_ms", float64(duration.Nanoseconds())/1e6),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

func panicMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := reqctx.GetRequestID(r.Context())

					logger.Error("panic recovered",
						slog.String("request_id", requestID),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("query", r.URL.RawQuery),
						slog.String("remote_addr", r.RemoteAddr),
						slog.String("user_agent", r.UserAgent()),
						slog.Any("panic", err),
						slog.String("stack_trace", string(debug.Stack())),
					)

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
