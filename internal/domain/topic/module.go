package topic

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
	group := r.Group("/topic", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/a", middleware.RequireRole("admin"))
	{
		admin.POST("/lesson/:lessonId", m.handler.CreateTopic)
		admin.GET("/lesson/:lessonId", m.handler.ListTopics)
		admin.PATCH("/:topicId", m.handler.UpdateTopic)
		admin.PATCH("/:topicId/reorder", m.handler.ReorderTopic)
		admin.DELETE("/:topicId", m.handler.DeleteTopic)
	}

	student := group.Group("/s", middleware.RequireRole("student"))
	{
		student.GET("/lesson/:lessonId", m.handler.ListMyTopics)
	}
}
