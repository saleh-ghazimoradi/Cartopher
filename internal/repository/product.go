package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryById(ctx context.Context, id uint) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id uint) error

	CreateProduct(ctx context.Context, product *domain.Product) error
	CreateProductImage(ctx context.Context, productImage *domain.ProductImage) error
	GetProductById(ctx context.Context, id uint) (*domain.Product, error)
	GetProducts(ctx context.Context, offset, limit int) ([]*domain.Product, error)
	GetProductImageCount(ctx context.Context, id uint) (int64, error)
	CountActiveProducts(ctx context.Context) (int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (p *productRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	return p.dbWrite.WithContext(ctx).Create(category).Error
}

func (p *productRepository) GetCategoryById(ctx context.Context, id uint) (*domain.Category, error) {
	var category *domain.Category
	if err := p.dbRead.WithContext(ctx).First(&category, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		}
		return nil, err
	}
	return category, nil
}

func (p *productRepository) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	if err := p.dbRead.WithContext(ctx).Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (p *productRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	return p.dbWrite.WithContext(ctx).Save(category).Error
}

func (p *productRepository) DeleteCategory(ctx context.Context, id uint) error {
	return p.dbWrite.WithContext(ctx).Delete(&domain.Category{}, id).Error
}

func (p *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	return p.dbWrite.WithContext(ctx).Create(product).Error
}

func (p *productRepository) CreateProductImage(ctx context.Context, productImage *domain.ProductImage) error {
	return p.dbWrite.WithContext(ctx).Create(productImage).Error
}

func (p *productRepository) GetProductById(ctx context.Context, id uint) (*domain.Product, error) {
	var product *domain.Product
	if err := p.dbRead.WithContext(ctx).Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		}
	}
	return product, nil
}

func (p *productRepository) GetProducts(ctx context.Context, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := p.dbRead.
		WithContext(ctx).
		Preload("Category").
		Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) GetProductImageCount(ctx context.Context, id uint) (int64, error) {
	var count int64
	if err := p.dbRead.WithContext(ctx).Model(&domain.ProductImage{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (p *productRepository) CountActiveProducts(ctx context.Context) (int64, error) {
	var total int64
	if err := p.dbRead.WithContext(ctx).Model(&domain.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return p.dbWrite.WithContext(ctx).Save(product).Error
}

func (p *productRepository) DeleteProduct(ctx context.Context, id uint) error {
	return p.dbWrite.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func NewProductRepository(dbWrite, dbRead *gorm.DB) ProductRepository {
	return &productRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
