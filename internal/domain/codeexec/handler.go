package codeexec

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type runRequest struct {
	AssignmentID string `json:"assignmentId,omitempty"`
	Source       string `json:"source"`
	Stdin        string `json:"stdin,omitempty"`
}

// Run executes a student's Python code once and returns stdout/stderr.
func (h *Handler) Run(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req runRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	res, err := h.svc.Run(c.Request.Context(), userID, req.Source, req.Stdin)
	if err != nil {
		writeRunError(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, res)
}

type gradePreviewRequest struct {
	Source string `json:"source"`
	Tests  []struct {
		Stdin          string `json:"stdin"`
		ExpectedStdout string `json:"expectedStdout"`
	} `json:"tests"`
}

type gradePreviewResponse struct {
	Passed  int             `json:"passed"`
	Total   int             `json:"total"`
	Results []PreviewResult `json:"results"`
}

func (h *Handler) GradePreview(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req gradePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	tests := make([]PreviewTest, len(req.Tests))
	for i, t := range req.Tests {
		tests[i] = PreviewTest{Stdin: t.Stdin, ExpectedStdout: t.ExpectedStdout}
	}

	passed, results, err := h.svc.GradePreview(c.Request.Context(), userID, req.Source, tests)
	if err != nil {
		writeRunError(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, gradePreviewResponse{
		Passed:  passed,
		Total:   len(results),
		Results: results,
	})
}

func writeRunError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		response.WriteError(c, http.StatusBadRequest, "no code to run")
	case errors.Is(err, ErrTooLarge):
		response.WriteError(c, http.StatusRequestEntityTooLarge, "code or input too large")
	case errors.Is(err, ErrRateLimited):
		response.WriteError(c, http.StatusTooManyRequests, "please wait a moment before running again")
	case errors.Is(err, ErrNotConfigured):
		response.WriteError(c, http.StatusServiceUnavailable, "code execution is not enabled")
	default:
		response.WriteInternal(c, err)
	}
}
