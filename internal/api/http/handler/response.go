package handler

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/assignment"

	"github.com/gin-gonic/gin"
)

func writeJSON(c *gin.Context, status int, data any) {
	if data == nil {
		return
	}
	c.JSON(status, data)
}

func writeError(c *gin.Context, status int, message string) {
	writeJSON(c, status, dto.ErrorResponse{Error: message})
}

func writeAssignmentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, assignment.ErrAccessDenied):
		writeError(c, http.StatusForbidden, err.Error())
	case errors.Is(err, assignment.ErrStudentNotFound):
		writeError(c, http.StatusNotFound, "student profile not found")
	default:
		writeError(c, http.StatusInternalServerError, err.Error())
	}
}