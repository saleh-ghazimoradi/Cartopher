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

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	order, err := o.orderService.CreateOrder(ctx, userId)
	if err != nil {
		helper.InternalServerError(ctx, "error while creating order", err)
		return
	}

	helper.CreatedResponse(ctx, "order successfully created", order)
}

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
