package service

import (
	"context"
	"errors"

	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
)

type CartService interface {
	GetCart(ctx context.Context, userId uint) (*dto.CartResponse, error)
	AddToCart(ctx context.Context, userId uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateCartItem(ctx context.Context, userId, itemId uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveFromCart(ctx context.Context, userId, itemId uint) error
}

type cartService struct {
	cartRepository    repository.CartRepository
	productRepository repository.ProductRepository
}

func (c *cartService) GetCart(ctx context.Context, userId uint) (*dto.CartResponse, error) {
	cart, err := c.cartRepository.GetCartByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return c.convertToCartResponse(cart), nil
}

func (c *cartService) AddToCart(ctx context.Context, userId uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	product, err := c.productRepository.GetProductById(ctx, req.ProductId)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock")
	}

	cart, err := c.cartRepository.GetOrCreateCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	cartItem, err := c.cartRepository.GetCartItem(ctx, cart.Id, req.ProductId)
	if err != nil {
		item := &domain.CartItem{
			CartId:    cart.Id,
			ProductId: product.Id,
			Quantity:  req.Quantity,
		}
		_ = c.cartRepository.CreateCartItem(ctx, item)
	} else {
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("not enough stock")
		}
		_ = c.cartRepository.UpdateCartItem(ctx, cartItem)
	}

	return c.GetCart(ctx, userId)
}

func (c *cartService) UpdateCartItem(ctx context.Context, userId, itemId uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cartItem, err := c.cartRepository.GetCartItemWithUser(ctx, userId, itemId)
	if err != nil {
		return nil, err
	}

	product, err := c.productRepository.GetProductById(ctx, cartItem.ProductId)
	if err != nil {
		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock")
	}

	cartItem.Quantity = req.Quantity
	if err := c.cartRepository.UpdateCartItem(ctx, cartItem); err != nil {
		return nil, err
	}

	return c.GetCart(ctx, userId)
}

func (c *cartService) RemoveFromCart(ctx context.Context, userId, itemId uint) error {
	return c.cartRepository.DeleteCartItem(ctx, userId, itemId)
}

func (c *cartService) convertToCartResponse(cart *domain.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			Id: cart.CartItems[i].Id,
			Product: dto.ProductResponse{
				Id:          cart.CartItems[i].Product.Id,
				CategoryId:  cart.CartItems[i].Product.CategoryId,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					Id:          cart.CartItems[i].Product.Category.Id,
					Name:        cart.CartItems[i].Product.Name,
					Description: cart.CartItems[i].Product.Description,
					IsActive:    cart.CartItems[i].Product.IsActive,
				},
			},
			Quantity: cart.CartItems[i].Quantity,
			Subtotal: subtotal,
		}
	}
	return &dto.CartResponse{
		Id:        cart.Id,
		UserId:    cart.UserId,
		CartItems: cartItems,
		Total:     total,
	}
}

func NewCartService(cartRepository repository.CartRepository, productRepository repository.ProductRepository) CartService {
	return &cartService{
		cartRepository:    cartRepository,
		productRepository: productRepository,
	}
}
