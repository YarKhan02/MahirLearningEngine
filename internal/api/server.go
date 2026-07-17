package api

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewServer(allowedOrigin string, addr string, modules []Module, logger *zap.Logger, rateLimitRequests int, rateLimitWindow time.Duration) *http.Server {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin, "https://www.mahircodelab.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RateLimit(logger, rateLimitRequests, rateLimitWindow))

	for _, m := range modules {
		m.RegisterRoutes(r)
	}
	
	return &http.Server{
		Addr: 				addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}