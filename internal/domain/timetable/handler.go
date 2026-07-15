package timetable

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/response"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTimetable(c *gin.Context) {
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req CreateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	t, err := ToCreateTimetable(batchID, req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Create(c.Request.Context(), t); err != nil {
		switch {
		case errors.Is(err, ErrNoWeekdays),
			errors.Is(err, ErrInvalidTime),
			errors.Is(err, ErrTimeOrder):
			response.WriteError(c, http.StatusBadRequest, err.Error())
		default:
			response.WriteError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.WriteJSON(c, http.StatusCreated, "successfully created timetable")
}

func (h *Handler) GetBatchTimetable(c *gin.Context) {
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	entries, err := h.svc.GetByBatch(c.Request.Context(), batchID)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]TimetableResponse, 0, len(entries))
	for _, e := range entries {
		resp = append(resp, ToTimetableResponse(e))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) DeleteTimetable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("timetableId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid timetable id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrTimetableNotFound) {
			response.WriteError(c, http.StatusNotFound, "timetable not found")
			return
		}
		response.WriteError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully deleted timetable")
}

func (h *Handler) GetMyUpcoming(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessions, err := h.svc.GetUpcomingForUser(c.Request.Context(), userID)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]ClassSessionResponse, 0, len(sessions))
	for _, s := range sessions {
		resp = append(resp, ToClassSessionResponse(s))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}
