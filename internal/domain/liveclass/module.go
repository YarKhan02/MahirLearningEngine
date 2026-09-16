package liveclass

import (
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/gin-gonic/gin"
)

// OriginPatterns returns the host[:port] patterns allowed to open the live-class
// WebSocket, matched against the browser Origin header (dev + production).
func OriginPatterns(allowedOrigin string) []string {
	return []string{
		stripScheme(allowedOrigin),
		"www.mahircodelab.com",
		"localhost:*",
		"127.0.0.1:*",
	}
}

func stripScheme(origin string) string {
	origin = strings.TrimPrefix(origin, "https://")
	origin = strings.TrimPrefix(origin, "http://")
	return strings.TrimSuffix(origin, "/")
}

type Module struct {
	handler  *Handler
	tokenSvc *token.Service
	redis    *redis.RedisClient
}

func NewModule(svc *Service, tokenSvc *token.Service, redis *redis.RedisClient, hub *Hub, originPatterns []string) *Module {
	return &Module{
		handler:  NewHandler(svc, hub, originPatterns),
		tokenSvc: tokenSvc,
		redis:    redis,
	}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/live/ws/:id", m.handler.ServeWS)

	authed := r.Group("/live", middleware.Auth(m.tokenSvc, m.redis))
	{
		authed.GET("/my/live", m.handler.GetMyLive)
		authed.GET("/my/batches", m.handler.MyBatches)
		authed.GET("/my/batch/:batchId/courses", m.handler.MyBatchCourses)
		authed.GET("/session/:id", m.handler.GetSession)
		authed.GET("/batch/:batchId/live", m.handler.GetLiveByBatch)
		authed.GET("/batch/:batchId/course/:courseId/sessions", m.handler.ListSessions)
		authed.POST("/:id/ticket", m.handler.IssueTicket)
		authed.POST("/:id/rtc-token", m.handler.RTCToken)
		authed.GET("/:id/whiteboards", m.handler.ListSnapshots)
	}

	admin := r.Group("/live/a", middleware.Auth(m.tokenSvc, m.redis), middleware.RequireRole("admin"))
	{
		admin.POST("/session", m.handler.StartSession)
		admin.POST("/session/:id/end", m.handler.EndSession)
		admin.POST("/:id/whiteboard/save", m.handler.SaveSnapshot)
		admin.POST("/:id/mic", m.handler.SetMic)
	}
}
