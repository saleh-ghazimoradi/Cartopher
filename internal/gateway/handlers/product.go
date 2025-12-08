package handlers

import "github.com/saleh-ghazimoradi/Cartopher/internal/service"

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}
