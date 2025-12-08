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

func NewProductRepository(dbWrite, dbRead *gorm.DB) ProductRepository {
	return &productRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
