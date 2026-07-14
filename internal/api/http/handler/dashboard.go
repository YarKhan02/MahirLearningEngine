package handler

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/dashboard"
	"github.com/gin-gonic/gin"
)

// DashboardHandler serves the admin landing-page overview.
type DashboardHandler struct {
	dashboardSvc *dashboard.Service
}

func NewDashboardHandler(dashboardSvc *dashboard.Service) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {

	d, err := h.dashboardSvc.GetAdminDashboard(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, mapper.ToAdminDashboardResponse(d))
}
