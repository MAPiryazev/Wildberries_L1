package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(rw, r)
			duration := time.Since(start)
			log.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

func JWTAuth(jwtSecret string, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("missing authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Warn("invalid authorization header format")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				log.Warn("failed to parse token", "err", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				log.Warn("invalid token claims")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				log.Warn("user_id not found in token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				log.Warn("invalid user_id format", "err", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			roleStr, ok := claims["role"].(string)
			if !ok {
				log.Warn("role not found in token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authCtx := &service.AuthContext{
				UserID: userID,
				Role:   models.Role(roleStr),
				Email:  claims["email"].(string),
			}

			ctx := service.SetUserContext(r.Context(), authCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
