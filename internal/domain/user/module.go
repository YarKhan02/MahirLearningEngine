package user

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/common"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler  	*Handler
	tokenSvc	*token.Service
	redis		*redis.RedisClient
}

func NewModule(svc *Service, studentSvc common.StudentProfileProvider, tokenSvc *token.Service, redis *redis.RedisClient, secureCookies bool) *Module {
	return &Module{
		handler: 	NewHandler(svc, studentSvc, tokenSvc, secureCookies),
		tokenSvc: 	tokenSvc,
		redis: 		redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/auth")
	
	public := group.Group("/public")
	{
		public.POST("/register", m.handler.RegisterAdmin)
		public.POST("/login", m.handler.Login)
		public.POST("/refresh", m.handler.Refresh)
		public.POST("/logout", m.handler.Logout)
	}
	authenticated := group.Group("/a", middleware.Auth(m.tokenSvc, m.redis))
	{
		authenticated.POST("/reset-password/:userID", m.handler.ResetPassword)
	}
}
