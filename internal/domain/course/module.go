package course

import "github.com/gin-gonic/gin"

type Module struct {
    handler *Handler
}

func NewModule(svc *Service) *Module {
    return &Module{handler: NewHandler(svc)}
}

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
    // same body as before
}