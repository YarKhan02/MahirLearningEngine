package course

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
    course := r.Group("/course", middleware.Auth(m.tokenSvc, m.redis), middleware.RequireRole("admin"))
	{
		course.GET("", m.handler.GetCourse)
		course.GET("/:courseId/lessons", m.handler.GetLesson)
		course.POST("", m.handler.InsertCourse)
		course.POST("/:courseId/lessons", m.handler.InsertLesson)
		course.PATCH("/lessons/:lessonId/reorder", m.handler.ReorderLesson)
		course.PATCH("/:courseId/lessons/:lessonId", m.handler.UpdateLesson)
		course.DELETE("/:courseId", m.handler.DeleteCourse)
	}
}