package handler

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/announcement"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnnouncementHandler struct {
	announcementSvc *announcement.Service
}

func NewAnnouncementHandler(announcementSvc *announcement.Service) *AnnouncementHandler {
	return &AnnouncementHandler{announcementSvc: announcementSvc}
}

func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	
	var req dto.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	a, err := mapper.ToCreateAnnouncement(req)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.announcementSvc.Create(c.Request.Context(), a); err != nil {
		switch {
		case errors.Is(err, announcement.ErrEmptyTitle), errors.Is(err, announcement.ErrEmptyDescription):
			writeError(c, http.StatusBadRequest, err.Error())
		default:
			writeError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(c, http.StatusCreated, "successfully posted announcement")
}

func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	list, err := h.announcementSvc.GetAll(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AnnouncementResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, mapper.ToAnnouncementResponse(a))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("announcementId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid announcement id")
		return
	}

	if err := h.announcementSvc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, announcement.ErrNotFound) {
			writeError(c, http.StatusNotFound, "announcement not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully deleted announcement")
}

func (h *AnnouncementHandler) GetMyAnnouncements(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.announcementSvc.GetForUser(c.Request.Context(), userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AnnouncementResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, mapper.ToAnnouncementResponse(a))
	}

	writeJSON(c, http.StatusOK, resp)
}
