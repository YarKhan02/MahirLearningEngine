package api

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(allowedOrigin string, addr string, modules []Module, logger *zap.Logger, rateLimitRequests int, rateLimitWindow time.Duration, PrometheusUsername string, PrometheusPassword string) *http.Server {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Recovery())
	r.Use(middleware.PrometheusMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin, "https://www.mahircodelab.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RateLimit(logger, rateLimitRequests, rateLimitWindow))

	metrics := r.Group(
		"/metrics",
		gin.BasicAuth(gin.Accounts{
			PrometheusUsername: PrometheusPassword,
		}),
	)
	metrics.GET("", gin.WrapH(promhttp.Handler()))

	for _, m := range modules {
		m.RegisterRoutes(r)
	}
	
	return &http.Server{
		Addr: 				addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}