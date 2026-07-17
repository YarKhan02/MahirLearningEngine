package announcement

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateAnnouncement(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	a, err := ToCreateAnnouncement(req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Create(c.Request.Context(), a); err != nil {
		switch {
		case errors.Is(err, ErrEmptyTitle), errors.Is(err, ErrEmptyDescription):
			response.WriteError(c, http.StatusBadRequest, err.Error())
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusCreated, "successfully posted announcement")
}

func (h *Handler) GetAnnouncements(c *gin.Context) {
	list, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AnnouncementResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, ToAnnouncementResponse(a))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) DeleteAnnouncement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("announcementId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid announcement id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "announcement not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully deleted announcement")
}

func (h *Handler) GetMyAnnouncements(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.svc.GetForUser(c.Request.Context(), userID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AnnouncementResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, ToAnnouncementResponse(a))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}
