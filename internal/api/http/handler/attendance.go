package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/attendance"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceHandler struct {
	attendanceSvc *attendance.Service
}

func NewAttendanceHandler(attendanceSvc *attendance.Service) *AttendanceHandler {
	return &AttendanceHandler{attendanceSvc: attendanceSvc}
}

// GetRoster lists a batch's students with their status for ?date= (default today).
func (h *AttendanceHandler) GetRoster(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format(constant.DateLayout)
	}

	date, err := time.Parse(constant.DateLayout, dateStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid date")
		return
	}

	roster, err := h.attendanceSvc.GetRoster(c.Request.Context(), batchIDU, date)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.RosterEntryResponse, 0, len(roster))
	for _, e := range roster {
		resp = append(resp, mapper.ToRosterEntryResponse(e))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AttendanceHandler) Mark(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req dto.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	date, err := time.Parse(constant.DateLayout, req.Date)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid date")
		return
	}

	studentIDU, err := uuid.Parse(req.StudentID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var createdBy uuid.UUID
	if userID, ok := middleware.CurrentUserID(c); ok {
		createdBy = userID
	}

	markAttendance := mapper.ToMarkAttendance(batchIDU, date, studentIDU, req.Status, createdBy)

	err = h.attendanceSvc.Mark(c.Request.Context(), markAttendance)
	if err != nil {
		if errors.Is(err, attendance.ErrInvalidStatus) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "attendance marked")
}

func (h *AttendanceHandler) GetStudentRecords(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	records, err := h.attendanceSvc.GetStudentRecords(c.Request.Context(), studentIDU)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AttendanceRecordResponse, 0, len(records))
	for _, rec := range records {
		resp = append(resp, mapper.ToAttendanceRecordResponse(rec))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AttendanceHandler) GetMyRecords(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	records, err := h.attendanceSvc.GetMyRecords(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, attendance.ErrStudentNotFound) {
			writeError(c, http.StatusNotFound, "student profile not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AttendanceRecordResponse, 0, len(records))
	for _, rec := range records {
		resp = append(resp, mapper.ToAttendanceRecordResponse(rec))
	}

	writeJSON(c, http.StatusOK, resp)
}
