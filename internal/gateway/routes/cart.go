package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type CartRoutes struct {
	cartHandler    *handlers.CartHandler
	authMiddleware *middlewares.Authentication
}

func (c *CartRoutes) cartRoute(router *gin.Engine) {
	v1 := router.Group("/v1")

	protected := v1.Group("/")
	protected.Use(c.authMiddleware.Authenticate())

	cart := protected.Group("/cart")
	cart.GET("/", c.cartHandler.GetCart)
	cart.POST("/items", c.cartHandler.AddToCart)
	cart.PUT("/items/:id", c.cartHandler.UpdateCart)
	cart.DELETE("/items/:id", c.cartHandler.RemoveCart)
}

func NewCartRoutes(cartHandler *handlers.CartHandler, authMiddleware *middlewares.Authentication) *CartRoutes {
	return &CartRoutes{
		cartHandler:    cartHandler,
		authMiddleware: authMiddleware,
	}
}
