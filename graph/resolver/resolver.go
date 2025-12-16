package resolver

import (
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"strconv"
)

type Resolver struct {
	authService    service.AuthService
	userService    service.UserService
	cartService    service.CartService
	orderService   service.OrderService
	productService service.ProductService
}

type Options func(*Resolver)

func WithAuthService(authService service.AuthService) Options {
	return func(r *Resolver) {
		r.authService = authService
	}
}

func WithUserService(userService service.UserService) Options {
	return func(r *Resolver) {
		r.userService = userService
	}
}

func WithCartService(cartService service.CartService) Options {
	return func(r *Resolver) {
		r.cartService = cartService
	}
}

func WithOrderService(orderService service.OrderService) Options {
	return func(r *Resolver) {
		r.orderService = orderService
	}
}

func WithProductService(productService service.ProductService) Options {
	return func(r *Resolver) {
		r.productService = productService
	}
}

func (r *Resolver) parseId(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	return uint(parsed), err
}

func NewResolver(options ...Options) *Resolver {
	resolver := &Resolver{}
	for _, option := range options {
		option(resolver)
	}

	return resolver
}
