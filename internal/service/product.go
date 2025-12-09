package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
)

type ProductService interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories(ctx context.Context) ([]*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error

	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	AddProductImage(ctx context.Context, productId uint, url, altText string) error
	GetProductById(ctx context.Context, id uint) (*dto.ProductResponse, error)
	GetProducts(ctx context.Context, page, limit int) ([]*dto.ProductResponse, *helper.PaginatedMeta, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func (p *productService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := p.productRepository.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (p *productService) GetCategories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := p.productRepository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoriesResponse := make([]*dto.CategoryResponse, len(categories))
	for i := range categories {
		categoriesResponse[i] = &dto.CategoryResponse{
			Id:          categories[i].Id,
			Name:        categories[i].Name,
			Description: categories[i].Description,
			IsActive:    categories[i].IsActive,
		}
	}

	return categoriesResponse, nil
}

func (p *productService) UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := p.productRepository.GetCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := p.productRepository.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (p *productService) DeleteCategory(ctx context.Context, id uint) error {
	return p.productRepository.DeleteCategory(ctx, id)
}

func (p *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &domain.Product{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	if err := p.productRepository.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return p.GetProductById(ctx, product.Id)

}

func (p *productService) AddProductImage(ctx context.Context, productId uint, url, altText string) error {
	count, err := p.productRepository.GetProductImageCount(ctx, productId)
	if err != nil {
		return err
	}

	image := &domain.ProductImage{
		ProductId: productId,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0,
	}

	return p.productRepository.CreateProductImage(ctx, image)
}

func (p *productService) GetProductById(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := p.productRepository.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	response := p.convertToProductResponse(product)
	return response, nil
}

func (p *productService) GetProducts(ctx context.Context, page, limit int) ([]*dto.ProductResponse, *helper.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	total, err := p.productRepository.CountActiveProducts(ctx)
	if err != nil {
		return nil, nil, err
	}

	products, err := p.productRepository.GetProducts(ctx, offset, limit)

	response := make([]*dto.ProductResponse, len(products))
	for i := range products {
		response[i] = p.convertToProductResponse(products[i])
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &helper.PaginatedMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPages,
	}

	return response, meta, err
}

func (p *productService) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := p.productRepository.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	product.CategoryId = req.CategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := p.productRepository.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}

	return p.GetProductById(ctx, product.Id)
}

func (p *productService) DeleteProduct(ctx context.Context, id uint) error {
	return p.productRepository.DeleteProduct(ctx, id)
}

func (p *productService) convertToProductResponse(product *domain.Product) *dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range images {
		images[i] = dto.ProductImageResponse{
			Id:        product.Images[i].Id,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
		}
	}

	return &dto.ProductResponse{
		Id:          product.Id,
		CategoryId:  product.CategoryId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			Id:          product.Category.Id,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
		},
		Images: images,
	}
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}
