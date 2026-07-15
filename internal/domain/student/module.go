package student

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type Module struct {
    handler 		*Handler
	userSvc 		*user.Service
	tempPassword	string
}

func NewModule(svc *Service, userSvc *user.Service, tempPassword string) *Module {
    return &Module{
		handler: NewHandler(svc, userSvc, tempPassword),
	}
}

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
    // same body as before
}