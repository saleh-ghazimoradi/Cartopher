package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type ProductRoutes struct {
	productHandler *handlers.ProductHandler
	authMiddleware *middlewares.Authentication
}

func (p *ProductRoutes) ProductRoute(router *gin.Engine) {
	v1 := router.Group("/v1")

	// Public routes
	v1.GET("/categories", p.productHandler.GetCategories)
	v1.GET("/products", p.productHandler.GetProducts)
	v1.GET("/products/:id", p.productHandler.GetProductById)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(p.authMiddleware.Authenticate())

	// Category admin routes
	categories := protected.Group("/categories")
	categories.POST("/", p.authMiddleware.AdminMiddleware(), p.productHandler.CreateCategory)
	categories.PUT("/:id", p.authMiddleware.AdminMiddleware(), p.productHandler.UpdateCategory)
	categories.DELETE("/:id", p.authMiddleware.AdminMiddleware(), p.productHandler.DeleteCategory)

	// Product admin routes
	products := protected.Group("/products")
	products.POST("/", p.authMiddleware.AdminMiddleware(), p.productHandler.CreateProduct)
	products.PUT("/:id", p.authMiddleware.AdminMiddleware(), p.productHandler.UpdateProduct)
	products.DELETE("/:id", p.authMiddleware.AdminMiddleware(), p.productHandler.DeleteProduct)
	products.POST("/:id/images", p.authMiddleware.AdminMiddleware(), p.productHandler.UploadProductImage)

}

func NewProductRoutes(productHandler *handlers.ProductHandler, authMiddleware *middlewares.Authentication) *ProductRoutes {
	return &ProductRoutes{
		productHandler: productHandler,
		authMiddleware: authMiddleware,
	}
}
