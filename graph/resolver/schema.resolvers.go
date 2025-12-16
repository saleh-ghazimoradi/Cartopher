package resolver

import (
	"context"
	"fmt"

	"github.com/saleh-ghazimoradi/Cartopher/graph"
	"github.com/saleh-ghazimoradi/Cartopher/graph/model"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input dto.RegisterRequest) (*dto.AuthResponse, error) {
	response, err := r.authService.Register(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("register failed: %w", err)
	}

	return response, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input dto.LoginRequest) (*dto.AuthResponse, error) {
	response, err := r.authService.Login(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	return response, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	response, err := r.authService.RefreshToken(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("refresh token failed: %w", err)
	}
	return response, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context, input dto.RefreshTokenRequest) (bool, error) {
	if err := r.authService.Logout(ctx, input.RefreshToken); err != nil {
		return false, fmt.Errorf("logout failed: %w", err)
	}
	return true, nil
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, input dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.userService.UpdateProfile(ctx, userId, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return user, nil
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, input dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	if !IsAdminFromContext(ctx) {
		return nil, ErrUnauthorized
	}

	category, err := r.productService.CreateCategory(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategory is the resolver for the updateCategory field.
func (r *mutationResolver) UpdateCategory(ctx context.Context, id string, input dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	if !IsAdminFromContext(ctx) {
		return nil, ErrUnauthorized
	}

	categoryId, err := r.parseId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse category id: %w", err)
	}

	category, err := r.productService.UpdateCategory(ctx, categoryId, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory is the resolver for the deleteCategory field.
func (r *mutationResolver) DeleteCategory(ctx context.Context, id string) (bool, error) {
	if !IsAdminFromContext(ctx) {
		return false, ErrUnauthorized
	}

	categoryId, err := r.parseId(id)
	if err != nil {
		return false, fmt.Errorf("failed to parse category id: %w", err)
	}

	if err := r.productService.DeleteCategory(ctx, categoryId); err != nil {
		return false, fmt.Errorf("failed to delete category: %w", err)
	}

	return true, nil
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if !IsAdminFromContext(ctx) {
		return nil, ErrUnauthorized
	}

	product, err := r.productService.CreateProduct(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if !IsAdminFromContext(ctx) {
		return nil, ErrUnauthorized
	}

	productId, err := r.parseId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse product id: %w", err)
	}

	product, err := r.productService.UpdateProduct(ctx, productId, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	if !IsAdminFromContext(ctx) {
		return false, ErrUnauthorized
	}

	productId, err := r.parseId(id)
	if err != nil {
		return false, fmt.Errorf("failed to parse product id: %w", err)
	}

	if err := r.productService.DeleteProduct(ctx, productId); err != nil {
		return false, fmt.Errorf("failed to delete product: %w", err)
	}

	return true, nil
}

// AddToCart is the resolver for the addToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, input dto.AddToCartRequest) (*dto.CartResponse, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	cart, err := r.cartService.AddToCart(ctx, userId, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	return cart, nil
}

// UpdateCartItem is the resolver for the updateCartItem field.
func (r *mutationResolver) UpdateCartItem(ctx context.Context, id string, input dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	itemId, err := r.parseId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse item id: %w", err)
	}

	cart, err := r.cartService.UpdateCartItem(ctx, userId, itemId, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return cart, nil
}

// RemoveFromCart is the resolver for the removeFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, id string) (bool, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return false, ErrUnauthorized
	}

	itemId, err := r.parseId(id)
	if err != nil {
		return false, fmt.Errorf("failed to parse item id: %w", err)
	}

	if err := r.cartService.RemoveFromCart(ctx, userId, itemId); err != nil {
		return false, fmt.Errorf("failed to remove from cart: %w", err)
	}

	return true, nil
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context) (*dto.OrderResponse, error) {
	panic(fmt.Errorf("not implemented: CreateOrder - createOrder"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*dto.UserResponse, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.userService.GetProfile(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	return user, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, page *int, limit *int) (*model.ProductConnection, error) {
	p := 1
	l := 10

	if page != nil {
		p = *page
	}

	if limit != nil {
		l = *limit
	}

	products, meta, err := r.productService.GetProducts(ctx, p, l)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	edges := make([]*model.ProductEdge, len(products))

	for i, product := range products {
		edges[i] = &model.ProductEdge{
			Node: product,
		}
	}

	return &model.ProductConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			Page:       meta.Page,
			Limit:      meta.Limit,
			Total:      int(meta.Total),
			TotalPages: meta.TotalPage,
		},
	}, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*dto.ProductResponse, error) {
	productId, err := r.parseId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse product id: %w", err)
	}

	product, err := r.productService.GetProductById(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	return product, nil
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := r.productService.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	result := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = category
	}

	return result, nil
}

// Cart is the resolver for the cart field.
func (r *queryResolver) Cart(ctx context.Context) (*dto.CartResponse, error) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	cart, err := r.cartService.GetCart(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}

	return cart, nil
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context, page *int, limit *int) (*model.OrderConnection, error) {
	panic(fmt.Errorf("not implemented: Orders - orders"))
}

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*dto.OrderResponse, error) {
	panic(fmt.Errorf("not implemented: Order - order"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
