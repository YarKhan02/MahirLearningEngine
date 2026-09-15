package program

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler  *Handler
	tokenSvc *token.Service
	redis    *redis.RedisClient
}

func NewModule(svc *Service, tokenSvc *token.Service, redis *redis.RedisClient) *Module {
	return &Module{
		handler:  NewHandler(svc),
		tokenSvc: tokenSvc,
		redis:    redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/program", m.handler.ListPublished)
	r.GET("/program/slug/:slug", m.handler.GetBySlug)

	admin := r.Group("/program/a", middleware.Auth(m.tokenSvc, m.redis), middleware.RequireRole("admin"))
	{
		admin.GET("", m.handler.ListAll)
		admin.POST("", m.handler.Create)
		admin.PATCH("/reorder", m.handler.Reorder)
		admin.PUT("/:id", m.handler.Update)
		admin.DELETE("/:id", m.handler.Delete)
		admin.PATCH("/:id/publish", m.handler.SetPublished)
	}
}
