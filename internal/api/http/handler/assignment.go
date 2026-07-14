package handler

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssignmentHandler struct {
	assignmentSvc *assignment.Service
}

func NewAssignmentHandler(assignmentSvc *assignment.Service) *AssignmentHandler {
	return &AssignmentHandler{assignmentSvc: assignmentSvc}
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	a, err := mapper.ToCreateAssignment(req, lessonIDU)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid due date")
		return
	}

	if err := h.assignmentSvc.CreateAssignment(c.Request.Context(), a); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, mapper.ToAssignmentResponse(*a))
}

func (h *AssignmentHandler) GetLessonAssignments(c *gin.Context) {

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	assignments, err := h.assignmentSvc.GetLessonAssignments(c.Request.Context(), lessonIDU)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		resp = append(resp, mapper.ToAssignmentResponse(a))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {

	assignmentIDU, err := uuid.Parse(c.Param("assignmentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	if err := h.assignmentSvc.DeleteAssignment(c.Request.Context(), assignmentIDU); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully deleted assignment")
}

func (h *AssignmentHandler) GetMyAssignments(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	assignments, err := h.assignmentSvc.GetStudentAssignments(c.Request.Context(), userID, lessonIDU)
	if err != nil {
		writeAssignmentError(c, err)
		return
	}

	resp := make([]dto.StudentAssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		resp = append(resp, mapper.ToStudentAssignmentResponse(a))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AssignmentHandler) SubmitAssignment(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	assignmentIDU, err := uuid.Parse(c.Param("assignmentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	var req dto.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.assignmentSvc.SubmitAssignment(c.Request.Context(), userID, assignmentIDU, req.Code); err != nil {
		writeAssignmentError(c, err)
		return
	}

	writeJSON(c, http.StatusOK, "submission saved")
}

func (h *AssignmentHandler) GetBatchSubmissions(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	submissions, err := h.assignmentSvc.GetBatchSubmissions(c.Request.Context(), batchIDU)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.BatchSubmissionResponse, 0, len(submissions))
	for _, s := range submissions {
		resp = append(resp, mapper.ToBatchSubmissionResponse(s))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *AssignmentHandler) GradeSubmission(c *gin.Context) {

	submissionIDU, err := uuid.Parse(c.Param("submissionId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid submission id")
		return
	}

	var req dto.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.assignmentSvc.GradeSubmission(c.Request.Context(), submissionIDU, *req.Marks, req.Remarks); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "submission graded")
}

func (h *AssignmentHandler) GetMySubmissions(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	submissions, err := h.assignmentSvc.GetMySubmissions(c.Request.Context(), userID)
	if err != nil {
		writeAssignmentError(c, err)
		return
	}

	resp := make([]dto.BatchSubmissionResponse, 0, len(submissions))
	for _, s := range submissions {
		resp = append(resp, mapper.ToBatchSubmissionResponse(s))
	}

	writeJSON(c, http.StatusOK, resp)
}
