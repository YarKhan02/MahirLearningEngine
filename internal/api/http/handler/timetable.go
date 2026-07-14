package handler

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/timetable"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TimetableHandler struct {
	timetableSvc *timetable.Service
}

func NewTimetableHandler(timetableSvc *timetable.Service) *TimetableHandler {
	return &TimetableHandler{timetableSvc: timetableSvc}
}

func (h *TimetableHandler) CreateTimetable(c *gin.Context) {
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req dto.CreateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	t, err := mapper.ToCreateTimetable(batchID, req)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.timetableSvc.Create(c.Request.Context(), t); err != nil {
		switch {
		case errors.Is(err, timetable.ErrNoWeekdays),
			errors.Is(err, timetable.ErrInvalidTime),
			errors.Is(err, timetable.ErrTimeOrder):
			writeError(c, http.StatusBadRequest, err.Error())
		default:
			writeError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(c, http.StatusCreated, "successfully created timetable")
}

func (h *TimetableHandler) GetBatchTimetable(c *gin.Context) {
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	entries, err := h.timetableSvc.GetByBatch(c.Request.Context(), batchID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.TimetableResponse, 0, len(entries))
	for _, e := range entries {
		resp = append(resp, mapper.ToTimetableResponse(e))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *TimetableHandler) DeleteTimetable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("timetableId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid timetable id")
		return
	}

	if err := h.timetableSvc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, timetable.ErrTimetableNotFound) {
			writeError(c, http.StatusNotFound, "timetable not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully deleted timetable")
}

func (h *TimetableHandler) GetMyUpcoming(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessions, err := h.timetableSvc.GetUpcomingForUser(c.Request.Context(), userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.ClassSessionResponse, 0, len(sessions))
	for _, s := range sessions {
		resp = append(resp, mapper.ToClassSessionResponse(s))
	}

	writeJSON(c, http.StatusOK, resp)
}
