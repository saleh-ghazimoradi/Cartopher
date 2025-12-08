package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
)

type ProductService interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories(ctx context.Context) ([]*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error
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

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}
