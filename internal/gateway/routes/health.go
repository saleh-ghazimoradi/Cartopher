package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
)

type HealthRoutes struct {
	HealthHandler *handlers.HealthHandler
}

func (h *HealthRoutes) HealthRoute(router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.GET("/health", h.HealthHandler.HealthCheck)
}

func NewHealthRoutes(healthRoutes *handlers.HealthHandler) *HealthRoutes {
	return &HealthRoutes{
		HealthHandler: healthRoutes,
	}
}
