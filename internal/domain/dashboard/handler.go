package dashboard

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAdminDashboard(c *gin.Context) {

	d, err := h.svc.GetAdminDashboard(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, ToAdminDashboardResponse(d))
}
