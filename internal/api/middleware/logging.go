package middleware

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger emits exactly one structured access-log line per request and
// seeds a request-scoped logger (request_id + client_ip) into the context so
// every downstream log line is correlated. CORS preflights are skipped as noise.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		requestID := c.GetString("request_id")
		reqLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("client_ip", c.ClientIP()),
		)
		ctx := logging.WithLogger(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		c.Next()

		// Route pattern (e.g. /batch/:batchId) keeps path cardinality low for Loki.
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		status := c.Writer.Status()
		fields := []zap.Field{
			zap.String("event", "http_request"),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Int64("latency_ms", time.Since(start).Milliseconds()),
		}

		if claims, ok := c.Get("claims"); ok {
			if cl, ok := claims.(*token.Claims); ok {
				fields = append(fields, zap.String("user_id", cl.UserID), zap.String("role", cl.Role))
			}
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		switch {
		case status >= 500:
			reqLogger.Error("request completed", fields...)
		case status >= 400:
			reqLogger.Warn("request completed", fields...)
		default:
			reqLogger.Info("request completed", fields...)
		}
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}
