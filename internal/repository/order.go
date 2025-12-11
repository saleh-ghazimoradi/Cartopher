package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByUserId(ctx context.Context, userId, orderId uint) (*domain.Order, error)
	GetOrderById(ctx context.Context, id uint) (*domain.Order, error)
	GetOrders(ctx context.Context, userId uint, offset, limit int) ([]domain.Order, error)
	CountOrders(ctx context.Context, userId uint) (int64, error)
	WithTx(tx *gorm.DB) OrderRepository
}

type orderRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
	tx      *gorm.DB
}

func (o *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	return exec(o.dbWrite, o.tx).WithContext(ctx).Create(&order).Error
}

func (o *orderRepository) GetOrderByUserId(ctx context.Context, userId, orderId uint) (*domain.Order, error) {
	var order domain.Order
	if err := exec(o.dbRead, o.tx).WithContext(ctx).Preload("OrderItems.Product.Category").Where("id = ? AND user_id = ?", orderId, userId).First(&order).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &order, nil
}

func (o *orderRepository) GetOrderById(ctx context.Context, id uint) (*domain.Order, error) {
	var order domain.Order
	if err := exec(o.dbRead, o.tx).WithContext(ctx).Preload("OrderItems.Product.Category").First(&order, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &order, nil
}

func (o *orderRepository) GetOrders(ctx context.Context, userId uint, offset, limit int) ([]domain.Order, error) {
	var orders []domain.Order
	if err := exec(o.dbRead, o.tx).WithContext(ctx).Preload("OrderItems.Product.Category").Where("user_id = ?", userId).Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return orders, nil
}

func (o *orderRepository) CountOrders(ctx context.Context, userId uint) (int64, error) {
	var count int64
	if err := exec(o.dbRead, o.tx).WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return 0, nil
		}
	}
	return count, nil
}

func (o *orderRepository) WithTx(tx *gorm.DB) OrderRepository {
	return &orderRepository{
		dbWrite: o.dbWrite,
		dbRead:  o.dbRead,
		tx:      tx,
	}
}

func NewOrderRepository(dbWrite, dbRead *gorm.DB) OrderRepository {
	return &orderRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
