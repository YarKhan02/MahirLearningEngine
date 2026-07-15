package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/handler"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerDeps struct {
    Config          *config.Config
    UserSvc         *user.Service
    RoleSvc         *role.Service
    CourseSvc       *course.Service
    BatchSvc        *batch.Service
    StudentSvc      *student.Service
    AssignmentSvc   *assignment.Service
    AttendanceSvc   *attendance.Service
    DashboardSvc    *dashboard.Service
    TimetableSvc    *timetable.Service
    AnnouncementSvc *announcement.Service
    TokenSvc        *token.Service
    Redis           *redis.RedisClient
}

func NewServer(deps ServerDeps) *http.Server {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{deps.Config.AllowedOrigin, "https://www.mahircodelab.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Secure cookies (SameSite=None) only work over HTTPS; in local HTTP dev we
	// must fall back to a plain SameSite=Lax cookie or the browser drops it.
	secureCookies := !strings.EqualFold(deps.Config.Env, "development")
	userHandler := handler.NewAuthHandler(deps.UserSvc, deps.StudentSvc, deps.TokenSvc, secureCookies)
	courseHandler := course.NewHandler(deps.CourseSvc)
	batchHandler := batch.NewHandler(deps.BatchSvc)
	studentHandler := student.NewHandler(deps.StudentSvc, deps.UserSvc, deps.Config.TempPassword)
	assignmentHandler := assignment.NewHandler(deps.AssignmentSvc)
	attendanceHandler := attendance.NewHandler(deps.AttendanceSvc)
	dashboardHandler := dashboard.NewHandler(deps.DashboardSvc)
	timetableHandler := timetable.NewHandler(deps.TimetableSvc)

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
	course := r.Group("/course", middleware.Auth(deps.TokenSvc, deps.Redis))
	
	// admin
	admin := course.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("", courseHandler.InsertCourse)
		admin.GET("", courseHandler.GetCourse)
		admin.DELETE("/:courseId", courseHandler.DeleteCourse)
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

	// student-facing course access (My Courses / lessons / progress)
	courseStudent := course.Group("/student", middleware.RequireRole("student"))
	{
		courseStudent.GET("", studentHandler.GetMyCourses)
		courseStudent.GET("/:courseId/lessons", studentHandler.GetMyLessons)
	}

	// student
	studentGroup := r.Group("/student", middleware.Auth(deps.TokenSvc, deps.Redis))
	
	// student/admin
	studentAdmin := studentGroup.Group("/admin", middleware.RequireRole("admin"),)
	{
		studentAdmin.GET("", studentHandler.GetStudents)
		studentAdmin.POST("", studentHandler.AdminCreateStudent)
		studentAdmin.PATCH("/:studentId/status", studentHandler.UpdateStudentStatus)
		studentAdmin.PATCH("/:studentId/batch", studentHandler.UpdateStudentBatch)
		studentAdmin.POST("/:studentId/account", studentHandler.CreateStudentAccount)
	}

	studentPortal := studentGroup.Group("/portal", middleware.RequireRole("student"))
	{
		studentPortal.POST("/lessons/:lessonId/progress", studentHandler.SetLessonProgress)
		studentPortal.GET("/lessons/:lessonId/assignments", assignmentHandler.GetMyAssignments)
		studentPortal.POST("/assignments/:assignmentId/submit", assignmentHandler.SubmitAssignment)
		studentPortal.GET("/submissions", assignmentHandler.GetMySubmissions)
		studentPortal.GET("/timetable", timetableHandler.GetMyUpcoming)
	}

	// batch/admin
	batch := r.Group("/batch", middleware.Auth(deps.TokenSvc, deps.Redis))
	admin = batch.Group("/admin")
	admin.Use(
		middleware.Auth(deps.TokenSvc, deps.Redis),
		middleware.RequireRole("admin"),
	)
	{
		admin.POST("", batchHandler.CreateBatch)
		admin.GET("", batchHandler.GetBatches)
		admin.PATCH("/:batchId", batchHandler.UpdateBatch)
		admin.DELETE("/:batchId", batchHandler.DeleteBatch)
		admin.GET("/:batchId/courses", batchHandler.GetBatchCourses)
		admin.PATCH("/:batchId/courses", batchHandler.UpdateBatchCourses)
		admin.POST("/:batchId/timetable", timetableHandler.CreateTimetable)
		admin.GET("/:batchId/timetable", timetableHandler.GetBatchTimetable)
		admin.DELETE("/timetable/:timetableId", timetableHandler.DeleteTimetable)
	}

	// attendance
	attendanceGroup := r.Group("/attendance", middleware.Auth(deps.TokenSvc, deps.Redis))
	
	// attendance/admin
	attendanceAdmin := attendanceGroup.Group("/admin", middleware.RequireRole("admin"))
	{
		attendanceAdmin.GET("/batches/:batchId", attendanceHandler.GetRoster)
		attendanceAdmin.POST("/batches/:batchId/mark", attendanceHandler.Mark)
		attendanceAdmin.GET("/students/:studentId", attendanceHandler.GetStudentRecords)
	}
	
	// attendance/student
	attendancePortal := attendanceGroup.Group("/portal", middleware.RequireRole("student"))
	{
		attendancePortal.GET("/me", attendanceHandler.GetMyRecords)
	}

	dashboardGroup := r.Group("/dashboard", middleware.Auth(deps.TokenSvc, deps.Redis))
	dashboardAdmin := dashboardGroup.Group("/admin", middleware.RequireRole("admin"))
	{
		dashboardAdmin.GET("", dashboardHandler.GetAdminDashboard)
	}

	return &http.Server{
		Addr: 				deps.Config.Addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}