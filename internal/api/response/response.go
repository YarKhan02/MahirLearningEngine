package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(c *gin.Context, status int, data any) {
	if data == nil {
		return
	}
	c.JSON(status, data)
}

func WriteError(c *gin.Context, status int, message string) {
	WriteJSON(c, status, ErrorResponse{Error: message})
}

func WriteInternal(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	WriteJSON(c, http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
}
