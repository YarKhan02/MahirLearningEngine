package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

func Auth(tokenSvc *token.Service, redis *redis.RedisClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := tokenSvc.ValidateAccessToken(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if redis != nil && claims.ID != "" {
				blocked, err := redis.Exists(r.Context(), "blocklist:" + claims.ID)
				if err != nil || blocked {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if claims.Role == role {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func GetClaims(ctx context.Context) *token.Claims {
	claims, _ := ctx.Value(claimsKey).(*token.Claims)
	return claims
}