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
	uploadService  service.UploadService
}

// CreateCategory docs
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} helper.Response{data=dto.CategoryResponse} "Category created successfully"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /categories [post]
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

// GetCategories docs
// @Summary Get all categories
// @Description Retrieve all active categories
// @Tags Categories
// @Produce json
// @Success 200 {object} helper.Response{data=[]dto.CategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} helper.Response "Internal server error"
// @Router /categories [get]
func (p *ProductHandler) GetCategories(ctx *gin.Context) {
	categories, err := p.productService.GetCategories(ctx)
	if err != nil {
		helper.InternalServerError(ctx, "Error getting categories", err)
		return
	}

	helper.SuccessResponse(ctx, "Categories successfully retrieved", categories)
}

// UpdateCategory docs
// @Summary Update a category
// @Description Update an existing category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} helper.Response{data=dto.CategoryResponse} "Category updated successfully"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /categories/{id} [put]
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

// DeleteCategory docs
// @Summary Delete a category
// @Description Delete a category (Admin only)
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} helper.Response "Category deleted successfully"
// @Failure 400 {object} helper.Response "Invalid category ID"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /categories/{id} [delete]
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

// CreateProduct docs
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} helper.Response{data=dto.ProductResponse} "Product created successfully"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /products [post]
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

// GetProducts
// @Summary Get all products
// @Description Retrieve paginated list of active products
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} helper.PaginatedResponse{data=[]dto.ProductResponse} "Products retrieved successfully"
// @Failure 500 {object} helper.Response "Internal server error"
// @Router /products [get]
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

// GetProductById docs
// @Summary Get a product by ID
// @Description Retrieve detailed information about a specific product
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} helper.Response{data=dto.ProductResponse} "Product retrieved successfully"
// @Failure 400 {object} helper.Response "Invalid product ID"
// @Failure 404 {object} helper.Response "Product not found"
// @Router /products/{id} [get]
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

// UpdateProduct docs
// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} helper.Response{data=dto.ProductResponse} "Product updated successfully"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /products/{id} [put]
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

// DeleteProduct docs
// @Summary Delete a product
// @Description Delete a product (Admin only)
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} helper.Response "Product deleted successfully"
// @Failure 400 {object} helper.Response "Invalid product ID"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /products/{id} [delete]
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

// UploadProductImage docs
// @Summary Upload product image
// @Description Upload an image for a product (Admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Image file"
// @Success 200 {object} helper.Response{data=map[string]string} "Image uploaded successfully"
// @Failure 400 {object} helper.Response "Invalid request or file"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 403 {object} helper.Response "Admin access required"
// @Router /products/{id}/images [post]
func (p *ProductHandler) UploadProductImage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "Invalid id given", err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		helper.BadRequestResponse(ctx, "No file uploaded", err)
		return
	}

	url, err := p.uploadService.UploadProductImage(uint(id), file)
	if err != nil {
		helper.InternalServerError(ctx, "Error uploading product image", err)
		return
	}

	if err := p.productService.AddProductImage(ctx, uint(id), url, file.Filename); err != nil {
		helper.InternalServerError(ctx, "Error adding product image", err)
		return
	}

	helper.SuccessResponse(ctx, "Product image successfully uploaded", url)
}

func NewProductHandler(productService service.ProductService, uploadService service.UploadService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		uploadService:  uploadService,
	}
}
