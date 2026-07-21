package quiz

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
	group := r.Group("/quiz", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/a", middleware.RequireRole("admin"))
	{
		admin.POST("/lesson/:lessonId", m.handler.CreateQuiz)
		admin.GET("/lesson/:lessonId", m.handler.ListQuizzes)
		admin.PUT("/:quizId", m.handler.EditQuiz)
		admin.DELETE("/:quizId", m.handler.DeleteQuiz)
		admin.GET("/:quizId/submissions", m.handler.ListSubmissions)
		admin.GET("/submission/:submissionId", m.handler.GetSubmission)
		admin.PATCH("/submission/:submissionId/grade", m.handler.GradeSubmission)
	}

	student := group.Group("/s", middleware.RequireRole("student"))
	{
		student.GET("/lesson/:lessonId", m.handler.ListMyQuizzes)
		student.POST("/:quizId/submit", m.handler.SubmitQuiz)
	}
}
