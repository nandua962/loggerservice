package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler is a handler responsible for handling health check requests.
func (genre *GenreController) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server run with base version",
	})
}
