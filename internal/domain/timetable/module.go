package timetable

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
		handler: NewHandler(svc),
		tokenSvc: 	tokenSvc,
		redis: 		redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
    group := r.Group("/timetable", middleware.Auth(m.tokenSvc, m.redis))
	admin := group.Group("/a", middleware.RequireRole("admin"))
    {

		admin.POST("/:batchId/timetable", m.handler.CreateTimetable)
		admin.GET("/:batchId/timetable", m.handler.GetBatchTimetable)
		admin.DELETE("/timetable/:timetableId", m.handler.DeleteTimetable)
	}
    student := group.Group("/s", middleware.RequireRole("student"))
    {
		student.GET("/timetable", m.handler.GetMyUpcoming)
	}
}