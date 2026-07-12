package http

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/handler"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, userSvc *user.Service, roleSvc *role.Service, courseSvc *course.Service, batchSvc *batch.Service, studentSvc *student.Service, assignmentSvc *assignment.Service, attendanceSvc *attendance.Service, tokenSvc *token.Service, redis *redis.RedisClient) *http.Server {
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
	studentHandler := handler.NewStudentHandler(studentSvc, userSvc, cfg.TempPassword)
	assignmentHandler := handler.NewAssignmentHandler(assignmentSvc)
	attendanceHandler := handler.NewAttendanceHandler(attendanceSvc)

	// Liveness probe for the deploy pipeline / Render health checks.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	user := r.Group("/auth")
	{
		user.POST("/register", userHandler.RegisterAdmin)
		user.POST("/login", userHandler.Login)
		user.POST("/refresh", userHandler.Refresh)
		user.POST("/logout", userHandler.Logout)
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
		admin.POST("/lessons/:lessonId/assignments", assignmentHandler.CreateAssignment)
		admin.GET("/lessons/:lessonId/assignments", assignmentHandler.GetLessonAssignments)
		admin.DELETE("/assignments/:assignmentId", assignmentHandler.DeleteAssignment)
		admin.GET("/batches/:batchId/submissions", assignmentHandler.GetBatchSubmissions)
		admin.PATCH("/submissions/:submissionId/grade", assignmentHandler.GradeSubmission)
	}

	// student
	studentGroup := r.Group("/student", middleware.Auth(tokenSvc, redis))
	
	// student/admin
	studentAdmin := studentGroup.Group("/admin", middleware.RequireRole("admin"),)
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
		studentPortal.GET("/lessons/:lessonId/assignments", assignmentHandler.GetMyAssignments)
		studentPortal.POST("/assignments/:assignmentId/submit", assignmentHandler.SubmitAssignment)
		studentPortal.GET("/submissions", assignmentHandler.GetMySubmissions)
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

	attendanceGroup := r.Group("/attendance", middleware.Auth(tokenSvc, redis))
	attendanceAdmin := attendanceGroup.Group("/admin", middleware.RequireRole("admin"))
	{
		attendanceAdmin.GET("/batches/:batchId", attendanceHandler.GetRoster)
		attendanceAdmin.POST("/batches/:batchId/mark", attendanceHandler.Mark)
		attendanceAdmin.GET("/students/:studentId", attendanceHandler.GetStudentRecords)
	}
	attendancePortal := attendanceGroup.Group("/portal", middleware.RequireRole("student"))
	{
		attendancePortal.GET("/me", attendanceHandler.GetMyRecords)
	}

	return &http.Server{
		Addr: 				cfg.Addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}