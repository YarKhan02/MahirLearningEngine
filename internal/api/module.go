package api

import "github.com/gin-gonic/gin"

type Module interface {
	RegisterRoutes(r *gin.Engine)
}