package http

import (
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/handler"
	"github.com/YarKhan02/MahirLearningEngine/internal/config"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/role"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, userSvc *user.Service, roleSvc *role.Service, courseSvc *course.Service, tokenSvc *token.Service) *http.Server {
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

	course := r.Group("course")
	{
		course.POST("admin", courseHandler.InsertCourse)
		course.GET("admin", courseHandler.GetCourse)
	}

	return &http.Server{
		Addr: 				cfg.Addr,
		Handler: 			r,
		ReadHeaderTimeout: 	5 * time.Second,
	}
}