package batch

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
        handler:    NewHandler(svc),
        tokenSvc: 	tokenSvc,
		redis: 		redis,
    }
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
    batch := r.Group("/batch", middleware.Auth(m.tokenSvc, m.redis), middleware.RequireRole("admin"))
	{
		batch.GET("/:batchId/courses", m.handler.GetBatchCourses)
		batch.GET("", m.handler.GetBatches)
		batch.POST("", m.handler.CreateBatch)
		batch.PATCH("/:batchId", m.handler.UpdateBatch)
		batch.PATCH("/:batchId/courses", m.handler.UpdateBatchCourses)
		batch.DELETE("/:batchId", m.handler.DeleteBatch)
    }
    public := r.Group("/public")
    {
        public.GET("/batches", m.handler.GetPublicBatches)
    }
}