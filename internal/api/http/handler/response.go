package handler

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	
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