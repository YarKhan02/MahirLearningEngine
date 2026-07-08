package http

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/handler"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, userSvc *user.Service, roleSvc *role.Service, courseSvc *course.Service, tokenSvc *token.Service, redis *redis.RedisClient) *http.Server {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	userHandler := handler.NewAuthHandler(userSvc, tokenSvc)
	courseHandler := handler.NewCourseHandler(courseSvc)

	user := r.Group("/auth")
	{
		user.POST("/register", userHandler.RegisterAdmin)
		user.POST("/login", userHandler.Login)
	}

	// /course/admin/
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

	return &http.Server{
		Addr: 				cfg.Addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}