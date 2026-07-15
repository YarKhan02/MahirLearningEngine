package announcement

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-gonic/gin"
)

type Module struct {
    handler     *Handler
    tokenSvc    *token.Service
    redis       *redis.RedisClient
}

func NewModule(svc *Service) *Module {
    return &Module{handler: NewHandler(svc)}
}

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
    group := r.Group("/announcement", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/admin", middleware.RequireRole("admin"))
	{
		admin.POST("", m.handler.CreateAnnouncement)
		admin.GET("", m.handler.GetAnnouncements)
		admin.DELETE("/:announcementId", m.handler.DeleteAnnouncement)
	}

	portal := group.Group("/portal", middleware.RequireRole("student"))
	{
		portal.GET("", m.handler.GetMyAnnouncements)
	}
}