package codeexec

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
	return &Module{handler: NewHandler(svc), tokenSvc: tokenSvc, redis: redis}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/code", middleware.Auth(m.tokenSvc, m.redis))
	{
		group.POST("/run", m.handler.Run)
	}
}
