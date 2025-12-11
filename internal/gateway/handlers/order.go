package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"strconv"
)

type OrderHandler struct {
	orderService service.OrderService
}

// CreateOrder dcs
// @Summary Create an order
// @Description Create an order from the current user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} helper.Response{data=dto.OrderResponse} "Order created successfully"
// @Failure 400 {object} helper.Response "Cart is empty or insufficient stock"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Router /orders [post]
func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	order, err := o.orderService.CreateOrder(ctx, userId)
	if err != nil {
		helper.InternalServerError(ctx, "error while creating order", err)
		return
	}

	helper.CreatedResponse(ctx, "order successfully created", order)
}

// GetOrder docs
// @Summary Get order by ID
// @Description Retrieve detailed information about a specific order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} helper.Response{data=dto.OrderResponse} "Order retrieved successfully"
// @Failure 400 {object} helper.Response "Invalid order ID"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 404 {object} helper.Response "Order not found"
// @Router /orders/{id} [get]
func (o *OrderHandler) GetOrder(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "invalid order id", err)
		return
	}

	order, err := o.orderService.GetOrder(ctx, userId, uint(id))
	if err != nil {
		helper.NotFoundResponse(ctx, "order not found")
		return
	}

	helper.SuccessResponse(ctx, "order successfully retrieved", order)
}

// GetOrders docs
// @Summary Get user's orders
// @Description Retrieve paginated list of user's orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} helper.PaginatedResponse{data=[]dto.OrderResponse} "Orders retrieved successfully"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 500 {object} helper.Response "Internal server error"
// @Router /orders [get]
func (o *OrderHandler) GetOrders(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	orders, meta, err := o.orderService.GetOrders(ctx, userId, page, limit)
	if err != nil {
		helper.InternalServerError(ctx, "error while getting orders", err)
		return
	}

	helper.PaginatedSuccessResponse(ctx, "Order retrieved successfully", orders, *meta)
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}
