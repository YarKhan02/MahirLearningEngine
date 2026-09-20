package maintenance

import "github.com/gin-gonic/gin"

type Module struct {
	handler *Handler
}

func NewModule(svc *Service, token string) *Module {
	return &Module{handler: NewHandler(svc, token)}
}

func (m *Module) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", m.handler.Health)
	r.POST("/internal/maintenance/cleanup", m.handler.Cleanup)
}
