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
	GetCartWithItemsAndProducts(ctx context.Context, userId uint) (*domain.Cart, error)
	GetOrCreateCart(ctx context.Context, userId uint) (*domain.Cart, error)
	GetCartItem(ctx context.Context, cartId, productId uint) (*domain.CartItem, error)
	CreateCartItem(ctx context.Context, item *domain.CartItem) error
	UpdateCartItem(ctx context.Context, item *domain.CartItem) error
	GetCartItemWithUser(ctx context.Context, userId, itemId uint) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, userId, itemId uint) error
	ClearCartItems(ctx context.Context, cartId uint) error
	WithTx(tx *gorm.DB) CartRepository
}

type cartRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
	tx      *gorm.DB
}

func (c *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	return exec(c.dbWrite, c.tx).WithContext(ctx).Create(cart).Error
}

func (c *cartRepository) GetCartByUserId(ctx context.Context, userId uint) (*domain.Cart, error) {
	var cart domain.Cart
	if err := exec(c.dbRead, c.tx).WithContext(ctx).Preload("CartItems.Product.Category").Where("user_id = ?", userId).First(&cart).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cart, nil
}

func (c *cartRepository) GetCartWithItemsAndProducts(ctx context.Context, userId uint) (*domain.Cart, error) {
	var cart domain.Cart
	if err := exec(c.dbRead, c.tx).WithContext(ctx).Preload("CartItems.Product").Where("user_id = ?", userId).First(&cart).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &cart, nil
}

func (c *cartRepository) GetOrCreateCart(ctx context.Context, userId uint) (*domain.Cart, error) {
	var cart domain.Cart

	err := exec(c.dbRead, c.tx).WithContext(ctx).
		Where("user_id = ?", userId).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = domain.Cart{UserId: userId}
			if err := exec(c.dbWrite, c.tx).WithContext(ctx).Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (c *cartRepository) GetCartItem(ctx context.Context, cartId, productId uint) (*domain.CartItem, error) {
	var cartItem domain.CartItem
	err := exec(c.dbRead, c.tx).WithContext(ctx).Where("cart_id = ? and product_id = ?", cartId, productId).First(&cartItem).Error
	return &cartItem, err
}

func (c *cartRepository) CreateCartItem(ctx context.Context, item *domain.CartItem) error {
	return exec(c.dbWrite, c.tx).WithContext(ctx).Create(item).Error
}

func (c *cartRepository) UpdateCartItem(ctx context.Context, item *domain.CartItem) error {
	return exec(c.dbWrite, c.tx).WithContext(ctx).Save(item).Error
}

func (c *cartRepository) GetCartItemWithUser(ctx context.Context, userId, itemId uint) (*domain.CartItem, error) {
	var item domain.CartItem
	if err := exec(c.dbRead, c.tx).WithContext(ctx).Joins("JOIN carts ON cart_items.cart_id = carts.id").
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
	return exec(c.dbWrite, c.tx).WithContext(ctx).Unscoped().Where(
		"id = ? AND cart_id IN (?)",
		itemId,
		c.dbRead.Select("id").Table("carts").Where("user_id = ?", userId),
	).Delete(&domain.CartItem{}).Error
}

func (c *cartRepository) ClearCartItems(ctx context.Context, cartId uint) error {
	return exec(c.dbWrite, c.tx).WithContext(ctx).Unscoped().Where("cart_id = ?", cartId).Delete(&domain.CartItem{}).Error
}

func (c *cartRepository) WithTx(tx *gorm.DB) CartRepository {
	return &cartRepository{
		dbWrite: c.dbWrite,
		dbRead:  c.dbRead,
		tx:      tx,
	}
}

func NewCartRepository(dbWrite, dbRead *gorm.DB) CartRepository {
	return &cartRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
