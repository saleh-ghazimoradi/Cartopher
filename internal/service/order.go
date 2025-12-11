package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userId uint) (*dto.OrderResponse, error)
	GetOrder(ctx context.Context, userId, orderId uint) (*dto.OrderResponse, error)
	GetOrders(ctx context.Context, userId uint, page, limit int) ([]*dto.OrderResponse, *helper.PaginatedMeta, error)
}

type orderService struct {
	orderRepository   repository.OrderRepository
	cartRepository    repository.CartRepository
	productRepository repository.ProductRepository
	db                *gorm.DB
}

func (o *orderService) CreateOrder(ctx context.Context, userId uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := o.db.Transaction(func(tx *gorm.DB) error {

		cartRepo := o.cartRepository.WithTx(tx)
		productRepo := o.productRepository.WithTx(tx)
		orderRepo := o.orderRepository.WithTx(tx)

		cart, err := cartRepo.GetCartWithItemsAndProducts(ctx, userId)
		if err != nil {
			return err
		}

		if len(cart.CartItems) == 0 {
			return errors.New("cart items is empty")
		}

		var totalAmount float64
		var orderItems []domain.OrderItem

		for i := range cart.CartItems {
			item := &cart.CartItems[i]

			if item.Product.Stock < item.Quantity {
				return fmt.Errorf("not enough stock for product: %s", item.Product.Name)
			}

			item.Product.Stock -= item.Quantity
			if err := productRepo.UpdateProduct(ctx, &item.Product); err != nil {
				return err
			}

			totalAmount += float64(item.Quantity) * item.Product.Price

			orderItems = append(orderItems, domain.OrderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			})
		}

		order := &domain.Order{
			UserId:      userId,
			Status:      domain.OrderStatusPending,
			TotalAmount: totalAmount,
			OrderItems:  orderItems,
		}

		if err := orderRepo.CreateOrder(ctx, order); err != nil {
			return err
		}

		if err := cartRepo.ClearCartItems(ctx, cart.Id); err != nil {
			return err
		}

		createdOrder, err := orderRepo.GetOrderById(ctx, order.Id)
		if err != nil {
			return err
		}

		resp := o.convertToOrderRepository(createdOrder)
		orderResponse = resp

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orderResponse, nil
}

func (o *orderService) GetOrders(ctx context.Context, userId uint, page, limit int) ([]*dto.OrderResponse, *helper.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	total, err := o.orderRepository.CountOrders(ctx, userId)
	if err != nil {
		return nil, nil, err
	}

	orders, err := o.orderRepository.GetOrders(ctx, userId, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	response := make([]*dto.OrderResponse, len(orders))
	for i := range orders {
		order := &orders[i]
		response[i] = o.convertToOrderRepository(order)
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	meta := &helper.PaginatedMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	return response, meta, nil
}

func (o *orderService) GetOrder(ctx context.Context, userId, orderId uint) (*dto.OrderResponse, error) {
	order, err := o.orderRepository.GetOrderByUserId(ctx, userId, orderId)
	if err != nil {
		return nil, err
	}

	response := o.convertToOrderRepository(order)
	return response, nil
}

func (o *orderService) getOrderResponse(ctx context.Context, orderId uint) (*dto.OrderResponse, error) {
	order, err := o.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	response := o.convertToOrderRepository(order)
	return response, nil
}

func (o *orderService) convertToOrderRepository(order *domain.Order) *dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))

	for i := range order.OrderItems {
		item := &order.OrderItems[i]
		orderItems[i] = dto.OrderItemResponse{
			Id: item.Id,
			Product: dto.ProductResponse{
				Id:          item.Product.Id,
				CategoryId:  item.Product.CategoryId,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				SKU:         item.Product.SKU,
				IsActive:    item.Product.IsActive,
				Category: dto.CategoryResponse{
					Id:          item.Product.Category.Id,
					Name:        item.Product.Category.Name,
					Description: item.Product.Category.Description,
					IsActive:    item.Product.Category.IsActive,
				},
			},
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		}
	}
	return &dto.OrderResponse{
		Id:          order.Id,
		UserId:      order.UserId,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt,
	}
}

func NewOrderService(orderRepository repository.OrderRepository, cartRepository repository.CartRepository, productRepository repository.ProductRepository, db *gorm.DB) OrderService {
	return &orderService{
		orderRepository:   orderRepository,
		cartRepository:    cartRepository,
		productRepository: productRepository,
		db:                db,
	}
}
