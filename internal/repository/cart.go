package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	CreateCart(ctx context.Context, cart *domain.Cart) error
}

type cartRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (r *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	return r.dbWrite.WithContext(ctx).Create(cart).Error
}

func NewCartRepository(dbWrite, dbRead *gorm.DB) CartRepository {
	return &cartRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
