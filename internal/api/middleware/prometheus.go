package middleware

import (
	"strconv"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Don't measure the scrape endpoint itself — it fires on every scrape
		// and would otherwise dominate the request counters.
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		c.Next()

		// Route template keeps cardinality flat; unmatched paths (scanners, 404s)
		// collapse into a single series instead of one per random URL.
		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

		if size := c.Writer.Size(); size > 0 {
			metrics.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(size))
		}
	}
}
