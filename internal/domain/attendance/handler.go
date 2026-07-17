package attendance

import (
	"errors"
	"net/http"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetRoster lists a batch's students with their status for ?date= (default today).
func (h *Handler) GetRoster(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format(constant.DateLayout)
	}

	date, err := time.Parse(constant.DateLayout, dateStr)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid date")
		return
	}

	roster, err := h.svc.GetRoster(c.Request.Context(), batchIDU, date)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]RosterEntryResponse, 0, len(roster))
	for _, e := range roster {
		resp = append(resp, ToRosterEntryResponse(e))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) Mark(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	date, err := time.Parse(constant.DateLayout, req.Date)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid date")
		return
	}

	studentIDU, err := uuid.Parse(req.StudentID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var createdBy uuid.UUID
	if userID, ok := middleware.CurrentUserID(c); ok {
		createdBy = userID
	}

	markAttendance := ToMarkAttendance(batchIDU, date, studentIDU, req.Status, createdBy)

	err = h.svc.Mark(c.Request.Context(), markAttendance)
	if err != nil {
		if errors.Is(err, ErrInvalidStatus) {
			response.WriteError(c, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "attendance marked")
}

func (h *Handler) GetStudentRecords(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	records, err := h.svc.GetStudentRecords(c.Request.Context(), studentIDU)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AttendanceRecordResponse, 0, len(records))
	for _, rec := range records {
		resp = append(resp, ToAttendanceRecordResponse(rec))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) GetMyRecords(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	records, err := h.svc.GetMyRecords(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student profile not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AttendanceRecordResponse, 0, len(records))
	for _, rec := range records {
		resp = append(resp, ToAttendanceRecordResponse(rec))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}
