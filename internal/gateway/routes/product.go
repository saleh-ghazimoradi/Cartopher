package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
)

type ProductRoutes struct {
	productHandler *handlers.ProductHandler
}

func (p *ProductRoutes) ProductRoute(router *gin.Engine) {}

func NewProductRoutes(productHandler *handlers.ProductHandler) *ProductRoutes {
	return &ProductRoutes{
		productHandler: productHandler,
	}
}
