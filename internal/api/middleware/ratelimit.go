package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*clientWindow
}

type clientWindow struct {
	start time.Time
	count int
}

// RateLimit is a fixed-window per-IP limiter. Breaches are logged as a
// structured security event so brute-force / flooding is visible in Grafana.
// A background janitor evicts stale windows so the map can't grow unbounded.
func RateLimit(logger *zap.Logger, limit int, window time.Duration) gin.HandlerFunc {
	rl := &rateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientWindow),
	}
	go rl.janitor()

	return func(c *gin.Context) {
		// Never throttle CORS preflight — it carries no credentials.
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if !rl.allow(ip) {
			metrics.RecordSecurityEvent("rate_limit_exceeded", "")
			logging.FromLogger(c.Request.Context()).Warn("rate limit exceeded",
				zap.String("event", "rate_limit_exceeded"),
				zap.String("path", c.FullPath()),
			)
			c.Header("Retry-After", window.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cw, ok := rl.clients[ip]
	if !ok || now.Sub(cw.start) >= rl.window {
		rl.clients[ip] = &clientWindow{start: now, count: 1}
		return true
	}

	if cw.count >= rl.limit {
		return false
	}

	cw.count++
	return true
}

// janitor periodically drops windows whose period has elapsed.
func (rl *rateLimiter) janitor() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for now := range ticker.C {
		rl.mu.Lock()
		for ip, cw := range rl.clients {
			if now.Sub(cw.start) >= rl.window {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
