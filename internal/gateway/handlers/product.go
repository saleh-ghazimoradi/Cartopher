package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"strconv"
)

type ProductHandler struct {
	productService service.ProductService
}

func (p *ProductHandler) CreateCategory(ctx *gin.Context) {
	var payload *dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid payload given", err)
		return
	}

	category, err := p.productService.CreateCategory(ctx, payload)
	if err != nil {
		helper.InternalServerError(ctx, "Error creating category", err)
		return
	}

	helper.CreatedResponse(ctx, "category successfully created", category)
}

func (p *ProductHandler) GetCategories(ctx *gin.Context) {
	categories, err := p.productService.GetCategories(ctx)
	if err != nil {
		helper.InternalServerError(ctx, "Error getting categories", err)
		return
	}

	helper.SuccessResponse(ctx, "Categories successfully retrieved", categories)
}

func (p *ProductHandler) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	var payload *dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid payload given", err)
		return
	}

	updatedCategory, err := p.productService.UpdateCategory(ctx, uint(id), payload)
	if err != nil {
		helper.InternalServerError(ctx, "Error updating category", err)
		return
	}

	helper.SuccessResponse(ctx, "Category successfully updated", updatedCategory)
}

func (p *ProductHandler) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	if err := p.productService.DeleteCategory(ctx, uint(id)); err != nil {
		helper.InternalServerError(ctx, "Error deleting category", err)
		return
	}

	helper.SuccessResponse(ctx, "Category successfully deleted", nil)
}

func (p *ProductHandler) CreateProduct(ctx *gin.Context) {
	var payload *dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid payload given", err)
		return
	}

	product, err := p.productService.CreateProduct(ctx, payload)
	if err != nil {
		helper.InternalServerError(ctx, "Error creating product", err)
		return
	}

	helper.CreatedResponse(ctx, "Product successfully created", product)
}

func (p *ProductHandler) GetProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, meta, err := p.productService.GetProducts(ctx, page, limit)
	if err != nil {
		helper.InternalServerError(ctx, "Error getting products", err)
		return
	}

	helper.PaginatedSuccessResponse(ctx, "Products retrieved successfully", products, *meta)
}

func (p *ProductHandler) GetProductById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	product, err := p.productService.GetProductById(ctx, uint(id))
	if err != nil {
		helper.NotFoundResponse(ctx, "Product not found")
		return
	}

	helper.SuccessResponse(ctx, "Product successfully retrieved", product)
}

func (p *ProductHandler) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	var payload *dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid payload given", err)
		return
	}

	updatedProduct, err := p.productService.UpdateProduct(ctx, uint(id), payload)
	if err != nil {
		helper.InternalServerError(ctx, "Error updating product", err)
		return
	}

	helper.SuccessResponse(ctx, "Product successfully updated", updatedProduct)
}

func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	if err := p.productService.DeleteProduct(ctx, uint(id)); err != nil {
		helper.InternalServerError(ctx, "Error deleting product", err)
		return
	}

	helper.SuccessResponse(ctx, "Product successfully deleted", nil)
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}
