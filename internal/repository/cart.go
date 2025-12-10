package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	CreateCart(ctx context.Context, cart *domain.Cart) error
	GetCartByUserId(ctx context.Context, userId uint) (*domain.Cart, error)
	GetOrCreateCart(ctx context.Context, userId uint) (*domain.Cart, error)
	GetCartItem(ctx context.Context, cartId, productId uint) (*domain.CartItem, error)
	CreateCartItem(ctx context.Context, item *domain.CartItem) error
	UpdateCartItem(ctx context.Context, item *domain.CartItem) error
	GetCartItemWithUser(ctx context.Context, userId, itemId uint) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, userId, itemId uint) error
}

type cartRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (c *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	return c.dbWrite.WithContext(ctx).Create(cart).Error
}

func (c *cartRepository) GetCartByUserId(ctx context.Context, userId uint) (*domain.Cart, error) {
	var cart domain.Cart
	if err := c.dbRead.WithContext(ctx).Preload("CartItems.Product.Category").Where("user_id = ?", userId).First(&cart).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cart, nil
}

func (c *cartRepository) GetOrCreateCart(ctx context.Context, userId uint) (*domain.Cart, error) {
	var cart domain.Cart

	err := c.dbRead.WithContext(ctx).
		Where("user_id = ?", userId).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = domain.Cart{UserId: userId}
			if err := c.dbWrite.WithContext(ctx).Create(&cart).Error; err != nil {
				return nil, err
			}
			// Return an empty cart
			return &cart, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (c *cartRepository) GetCartItem(ctx context.Context, cartId, productId uint) (*domain.CartItem, error) {
	var cartItem domain.CartItem
	err := c.dbRead.WithContext(ctx).Where("cart_id = ? and product_id = ?", cartId, productId).First(&cartItem).Error
	return &cartItem, err
}

func (c *cartRepository) CreateCartItem(ctx context.Context, item *domain.CartItem) error {
	return c.dbWrite.WithContext(ctx).Create(item).Error
}

func (c *cartRepository) UpdateCartItem(ctx context.Context, item *domain.CartItem) error {
	return c.dbWrite.WithContext(ctx).Save(item).Error
}

func (c *cartRepository) GetCartItemWithUser(ctx context.Context, userId, itemId uint) (*domain.CartItem, error) {
	var item domain.CartItem
	if err := c.dbRead.WithContext(ctx).Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemId, userId).
		First(&item).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

func (c *cartRepository) DeleteCartItem(ctx context.Context, userId, itemId uint) error {
	return c.dbWrite.WithContext(ctx).Where(
		"id = ? AND cart_id IN (?)",
		itemId,
		c.dbRead.Select("id").Table("carts").Where("user_id = ?", userId),
	).Delete(&domain.CartItem{}).Error
}

func NewCartRepository(dbWrite, dbRead *gorm.DB) CartRepository {
	return &cartRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
