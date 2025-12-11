package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type OrderRoutes struct {
	orderHandler   *handlers.OrderHandler
	authMiddleware *middlewares.Authentication
}

func (o *OrderRoutes) OrderRoute(router *gin.Engine) {
	v1 := router.Group("/v1")
	protected := v1.Group("/")
	protected.Use(o.authMiddleware.Authenticate())
	orders := protected.Group("/orders")
	orders.POST("/", o.orderHandler.CreateOrder)
	orders.GET("/", o.orderHandler.GetOrders)
	orders.GET("/:id", o.orderHandler.GetOrder)
}

func NewOrderRoutes(orderHandler *handlers.OrderHandler, authMiddleware *middlewares.Authentication) *OrderRoutes {
	return &OrderRoutes{
		orderHandler:   orderHandler,
		authMiddleware: authMiddleware,
	}
}
