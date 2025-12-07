package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "I'm breathing!"})
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
