package assignment

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func submissionStatusFilter(c *gin.Context) string {
	switch s := c.Query("status"); s {
	case "submitted", "graded":
		return s
	default:
		return ""
	}
}

func (h *Handler) CreateAssignment(c *gin.Context) {

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	a, err := ToCreateAssignment(req, lessonIDU)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid due date")
		return
	}

	if err := h.svc.CreateAssignment(c.Request.Context(), a); err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, ToAssignmentResponse(*a))
}

func (h *Handler) GetLessonAssignments(c *gin.Context) {

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	assignments, err := h.svc.GetLessonAssignments(c.Request.Context(), lessonIDU)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		resp = append(resp, ToAssignmentResponse(a))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) DeleteAssignment(c *gin.Context) {

	assignmentIDU, err := uuid.Parse(c.Param("assignmentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	if err := h.svc.DeleteAssignment(c.Request.Context(), assignmentIDU); err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully deleted assignment")
}

func (h *Handler) GetMyAssignments(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	assignments, err := h.svc.GetStudentAssignments(c.Request.Context(), userID, lessonIDU)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := make([]StudentAssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		resp = append(resp, ToStudentAssignmentResponse(a))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) SubmitAssignment(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	assignmentIDU, err := uuid.Parse(c.Param("assignmentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	var req SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.SubmitAssignment(c.Request.Context(), userID, assignmentIDU, req.Code); err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(c, http.StatusOK, "submission saved")
}

func (h *Handler) GetBatchSubmissions(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	p := pagination.Parse(c.Query("page"), c.Query("pageSize"), 20, 100)

	submissions, total, err := h.svc.GetBatchSubmissions(
		c.Request.Context(), batchIDU, c.Query("q"), submissionStatusFilter(c), p.Limit(), p.Offset(),
	)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]BatchSubmissionResponse, 0, len(submissions))
	for _, s := range submissions {
		resp = append(resp, ToBatchSubmissionResponse(s))
	}

	response.WriteJSON(c, http.StatusOK, pagination.NewPage(resp, total, p))
}

func (h *Handler) GetBatchSubmissionSummary(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	summary, err := h.svc.GetBatchSubmissionSummary(c.Request.Context(), batchIDU, c.Query("q"))
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, ToSubmissionSummaryResponse(summary))
}

func (h *Handler) GradeSubmission(c *gin.Context) {

	submissionIDU, err := uuid.Parse(c.Param("submissionId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid submission id")
		return
	}

	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.GradeSubmission(c.Request.Context(), submissionIDU, *req.Marks, req.Remarks); err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "submission graded")
}

func (h *Handler) GetMySubmissions(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	p := pagination.Parse(c.Query("page"), c.Query("pageSize"), 20, 100)

	submissions, total, err := h.svc.GetMySubmissions(
		c.Request.Context(), userID, submissionStatusFilter(c), p.Limit(), p.Offset(),
	)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student profile not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	resp := make([]BatchSubmissionResponse, 0, len(submissions))
	for _, s := range submissions {
		resp = append(resp, ToBatchSubmissionResponse(s))
	}

	response.WriteJSON(c, http.StatusOK, pagination.NewPage(resp, total, p))
}

// GetMySubmissionSummary returns the logged-in student's per-status counts.
func (h *Handler) GetMySubmissionSummary(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.svc.GetMySubmissionSummary(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student profile not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, ToSubmissionSummaryResponse(summary))
}
