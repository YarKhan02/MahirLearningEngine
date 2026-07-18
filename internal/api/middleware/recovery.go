package middleware

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery converts a panic into a structured error log + a generic 500,
// so a single bad request can never take the process down or leak internals.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		metrics.RecordPanic()
		logging.FromLogger(c.Request.Context()).Error("panic recovered",
			zap.String("event", "panic"),
			zap.Any("error", err),
			zap.String("path", c.FullPath()),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
	})
}
