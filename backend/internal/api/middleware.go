package api

import (
	auth "astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/repository"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func AuthMiddleware(sessions *repository.SessionRepository, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				problem(w, http.StatusUnauthorized, "missing token")
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				problem(w, http.StatusUnauthorized, "invalid token")
				return
			}
			claims, err := auth.ValidateToken(parts[1], secret)
			if err != nil {
				problem(w, http.StatusUnauthorized, "invalid token")
				return
			}
			jti, _ := claims["jti"].(string)
			userID, _ := claims["user_id"].(string)
			role, _ := claims["role"].(string)
			_, err = sessions.Get(r.Context(), jti)
			if err != nil {
				problem(w, http.StatusUnauthorized, "session not found")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "role", role)
			ctx = context.WithValue(ctx, "jti", jti)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
