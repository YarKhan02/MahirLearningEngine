package http

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/handler"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, userSvc *user.Service, roleSvc *role.Service, courseSvc *course.Service, batchSvc *batch.Service, studentSvc *student.Service, tokenSvc *token.Service, redis *redis.RedisClient) *http.Server {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	userHandler := handler.NewAuthHandler(userSvc, tokenSvc)
	courseHandler := handler.NewCourseHandler(courseSvc)
	batchHandler := handler.NewBatchHandler(batchSvc)
	studentHandler := handler.NewStudentHandler(studentSvc, userSvc)

	user := r.Group("/auth")
	{
		user.POST("/register", userHandler.RegisterAdmin)
		user.POST("/login", userHandler.Login)
	}

	// public routes
	public := r.Group("/public")
	{
		public.GET("/batches", batchHandler.GetPublicBatches)
		public.POST("/students/register", studentHandler.RegisterStudent)
	}

	// course/admin/
	course := r.Group("/course", middleware.Auth(tokenSvc, redis))
	admin := course.Group("/admin")
	admin.Use(
		middleware.Auth(tokenSvc, redis),
		middleware.RequireRole("admin"),
	)
	{
		admin.POST("", courseHandler.InsertCourse)
		admin.GET("", courseHandler.GetCourse)
		admin.POST("/:courseId/lessons", courseHandler.InsertLesson)
		admin.GET("/:courseId/lessons", courseHandler.GetLesson)
		admin.PATCH("/:courseId/lessons/:lessonId", courseHandler.UpdateLesson)
		admin.PATCH("/lessons/:lessonId/reorder", courseHandler.ReorderLesson)
	}

	// student/admin/
	studentGroup := r.Group("/student", middleware.Auth(tokenSvc, redis))
	studentAdmin := studentGroup.Group("/admin")
	studentAdmin.Use(
		middleware.Auth(tokenSvc, redis),
		middleware.RequireRole("admin"),
	)
	{
		studentAdmin.GET("", studentHandler.GetStudents)
		studentAdmin.POST("", studentHandler.AdminCreateStudent)
		studentAdmin.PATCH("/:studentId/status", studentHandler.UpdateStudentStatus)
		studentAdmin.PATCH("/:studentId/batch", studentHandler.UpdateStudentBatch)
		studentAdmin.POST("/:studentId/account", studentHandler.CreateStudentAccount)
	}

	// student-facing course access (My Courses / lessons / progress)
	courseStudent := course.Group("/student", middleware.RequireRole("student"))
	{
		courseStudent.GET("", studentHandler.GetMyCourses)
		courseStudent.GET("/:courseId/lessons", studentHandler.GetMyLessons)
	}

	studentPortal := studentGroup.Group("/portal", middleware.RequireRole("student"))
	{
		studentPortal.POST("/lessons/:lessonId/progress", studentHandler.SetLessonProgress)
	}

	// batch/admin
	batch := r.Group("/batch", middleware.Auth(tokenSvc, redis))
	admin = batch.Group("/admin")
	admin.Use(
		middleware.Auth(tokenSvc, redis),
		middleware.RequireRole("admin"),
	)
	{
		admin.POST("", batchHandler.CreateBatch)
		admin.GET("", batchHandler.GetBatches)
		admin.GET("/:batchId/courses", batchHandler.GetBatchCourses)
		admin.PATCH("/:batchId/courses", batchHandler.UpdateBatchCourses)
	}

	return &http.Server{
		Addr: 				cfg.Addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}