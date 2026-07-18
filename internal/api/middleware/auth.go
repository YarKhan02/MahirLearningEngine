package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// logAuthFailure records a rejected request as a structured security event.
// client_ip and request_id are already on the context logger.
func logAuthFailure(c *gin.Context, reason string) {
	metrics.RecordSecurityEvent("auth_failure", reason)
	logging.FromLogger(c.Request.Context()).Warn("authentication failed",
		zap.String("event", "auth_failure"),
		zap.String("reason", reason),
		zap.String("path", c.FullPath()),
	)
}

func Auth(tokenSvc *token.Service, redis *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logAuthFailure(c, "missing_header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logAuthFailure(c, "malformed_header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := tokenSvc.ValidateAccessToken(parts[1])
		if err != nil {
			logAuthFailure(c, "invalid_token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if redis != nil && claims.ID != "" {
			blocked, err := redis.Exists(c.Request.Context(), "blocklist:" + claims.ID)
			if err != nil || blocked {
				logAuthFailure(c, "token_revoked")
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
			metrics.RecordSecurityEvent("forbidden", "role_mismatch")
			logging.FromLogger(c.Request.Context()).Warn("access forbidden",
				zap.String("event", "forbidden"),
				zap.String("required_role", role),
				zap.String("actual_role", claims.Role),
				zap.String("user_id", claims.UserID),
				zap.String("path", c.FullPath()),
			)
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
// CurrentUser returns the authenticated user's claims set by Auth.
func CurrentUser(c *gin.Context) (*token.Claims, bool) {
	value, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claims, ok := value.(*token.Claims)
	return claims, ok
}

// CurrentUserID returns the authenticated user's id from the JWT claims.
func CurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	claims, ok := CurrentUser(c)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := claims.UserUUID()
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}
