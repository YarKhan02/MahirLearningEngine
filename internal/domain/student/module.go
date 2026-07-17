package student

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/common"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler  *Handler
	tokenSvc *token.Service
	redis    *redis.RedisClient
}

func NewModule(svc *Service, userSvc common.StudentAccountRegistrar, tokenSvc *token.Service, redis *redis.RedisClient, tempPassword string) *Module {
	return &Module{
		handler:  	NewHandler(svc, userSvc, tempPassword),
		tokenSvc: 	tokenSvc,
		redis:		redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	// public route
	public := r.Group("/public")
	{
		public.POST("/students/register", m.handler.RegisterStudent)
	}

	group := r.Group("/dashboard", middleware.Auth(m.tokenSvc, m.redis))

	admin := group.Group("/admin", middleware.RequireRole("admin"))
	{
		admin.GET("", m.handler.GetStudents)
		admin.POST("", m.handler.AdminCreateStudent)
		admin.PATCH("/:studentId/status", m.handler.UpdateStudentStatus)
		admin.PATCH("/:studentId/batch", m.handler.UpdateStudentBatch)
		admin.POST("/:studentId/account", m.handler.CreateStudentAccount)
	}
	student := group.Group("/student", middleware.RequireRole("student"))
	{
		student.GET("", m.handler.GetMyCourses)
		student.GET("/:courseId/lessons", m.handler.GetMyLessons)
		student.POST("/:lessonId/progress", m.handler.SetLessonProgress)
	}
}
