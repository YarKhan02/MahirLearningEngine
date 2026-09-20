package maintenance

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc   *Service
	token string // shared secret for the cron trigger; empty = endpoint disabled
}

func NewHandler(svc *Service, token string) *Handler {
	return &Handler{svc: svc, token: token}
}

// Health is a lightweight readiness probe (DB ping) for Render + uptime pings.
func (h *Handler) Health(c *gin.Context) {
	if err := h.svc.Ping(c.Request.Context()); err != nil {
		response.WriteError(c, http.StatusServiceUnavailable, "unhealthy")
		return
	}
	response.WriteJSON(c, http.StatusOK, gin.H{"status": "ok"})
}

// Cleanup runs the maintenance jobs
func (h *Handler) Cleanup(c *gin.Context) {
	if h.token == "" {
		response.WriteError(c, http.StatusServiceUnavailable, "maintenance is not configured")
		return
	}
	if !h.authorized(c) {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	report := h.svc.RunAll(c.Request.Context())
	response.WriteJSON(c, http.StatusOK, report)
}

// authorized accepts the secret via `Authorization: Bearer <token>` or
// `X-Maintenance-Token: <token>`, compared in constant time.
func (h *Handler) authorized(c *gin.Context) bool {
	provided := c.GetHeader("X-Maintenance-Token")
	if provided == "" {
		if after, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer "); ok {
			provided = after
		}
	}
	if provided == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(h.token)) == 1
}
