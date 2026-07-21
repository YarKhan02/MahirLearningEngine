package attachement

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

type Module struct {
    handler     *Handler
    tokenSvc    *token.Service
    redis       *redis.RedisClient
}

func NewModule(svc *Service, tokenSvc *token.Service, redis *redis.RedisClient) *Module {
    return &Module{
		handler: 	NewHandler(svc),
		tokenSvc: 	tokenSvc,
		redis: 		redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/attachment", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/a", middleware.RequireRole("admin"))
	{
		admin.POST("/presign", m.handler.PresignUpload)
		admin.POST("/confirm", m.handler.ConfirmUpload)
		admin.GET("/course/:courseId", m.handler.ListCourseMaterials)
		admin.DELETE("/:attachmentId", m.handler.DeleteMaterial)
	}

	student := group.Group("/s", middleware.RequireRole("student"))
	{
		student.GET("/course/:courseId", m.handler.ListMyCourseMaterials)
	}
}