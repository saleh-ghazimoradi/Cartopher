package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the server is running
// @Tags         health
// @Produce      json
// @Success      200 {string} string "I'm breathing!"
// @Router       /health [get]
func (h *HealthHandler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "I'm breathing!"})
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
