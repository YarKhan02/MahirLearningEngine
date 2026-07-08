package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

func Auth(tokenSvc *token.Service, redis *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := tokenSvc.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if redis != nil && claims.ID != "" {
			blocked, err := redis.Exists(c.Request.Context(), "blocklist:" + claims.ID)
			if err != nil || blocked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token revoked",
				})
				return
			}
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		
		claims, ok := value.(*token.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid claims",
			})
			return
		}
		
		if claims.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}
		c.Next()
	}
}

func GetClaims(ctx context.Context) *token.Claims {
	claims, _ := ctx.Value(claimsKey).(*token.Claims)
	return claims
}