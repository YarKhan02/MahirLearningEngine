package response

import (
	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"

	"github.com/gin-gonic/gin"
)

func WriteJSON(c *gin.Context, status int, data any) {
	if data == nil {
		return
	}
	c.JSON(status, data)
}

func WriteError(c *gin.Context, status int, message string) {
	WriteJSON(c, status, dto.ErrorResponse{Error: message})
}