package assignment

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
    group := r.Group("/assignment", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/a", middleware.RequireRole("admin"))
    {
		admin.GET("/:lessonId/assignments", m.handler.GetLessonAssignments)
		admin.GET("/batch/:batchId/submissions", m.handler.GetBatchSubmissions)
        admin.POST("/:lessonId/assignments", m.handler.CreateAssignment)
		admin.PATCH("/:submissionId/grade", m.handler.GradeSubmission)
		admin.DELETE("/:assignmentId", m.handler.DeleteAssignment)
    }
    student := group.Group("/s", middleware.RequireRole("student"))
    {
        student.GET("/lessons/:lessonId/assignments", m.handler.GetMyAssignments)
		student.GET("/submissions", m.handler.GetMySubmissions)
		student.POST("/assignments/:assignmentId/submit", m.handler.SubmitAssignment)
    }
}