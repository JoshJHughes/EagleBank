package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type ServerArgs struct {
	Logger  *slog.Logger
	UserSvc UserService
	AcctSvc AccountService
	TanSvc  TransactionService
}

func NewServer(args ServerArgs) http.Handler {
	mux := http.NewServeMux()

	// unprotected routes
	mux.HandleFunc("/health", handleHealth())
	mux.HandleFunc("POST /login", handleLogin())
	mux.HandleFunc("POST /v1/users", handleCreateUser(args.UserSvc))

	// protected routes
	mux.HandleFunc("GET /v1/users/{userId}", authMiddleware(handleGetUser(args.UserSvc)))

	mux.HandleFunc("POST /v1/accounts", authMiddleware(handleCreateAccount(args.AcctSvc)))
	mux.HandleFunc("GET /v1/accounts", authMiddleware(handleListAccounts(args.AcctSvc)))
	mux.HandleFunc("GET /v1/accounts/{accountId}", authMiddleware(handleFetchAccount(args.AcctSvc)))

	mux.HandleFunc("POST /v1/accounts/{accountId}/transactions", authMiddleware(handleCreateTransaction(args.TanSvc, args.AcctSvc)))

	handler := panicMiddleware(args.Logger)(mux)
	handler = loggingMiddleware(args.Logger)(handler)

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

			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
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
					requestID := GetRequestID(r.Context())

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

func authMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("authorization header required"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("invalid authorization header format"))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("invalid token claims"))
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok {
			writeErrorResponse(w, http.StatusUnauthorized, errors.New("missing userID in token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
