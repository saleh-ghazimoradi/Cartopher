package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"strconv"
)

type CartHandler struct {
	cartService service.CartService
}

func (c *CartHandler) GetCart(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	cart, err := c.cartService.GetCart(ctx, userId)
	if err != nil {
		helper.NotFoundResponse(ctx, "Cart not found")
		return
	}

	helper.SuccessResponse(ctx, "Cart retrieved successfully", cart)
}

func (c *CartHandler) AddToCart(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	var payload dto.AddToCartRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "invalid payload given", err)
		return
	}

	cart, err := c.cartService.AddToCart(ctx, userId, &payload)
	if err != nil {
		helper.InternalServerError(ctx, "error adding to cart", err)
		return
	}

	helper.SuccessResponse(ctx, "item added to cart successfully", cart)
}

func (c *CartHandler) UpdateCart(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "invalid id", err)
		return
	}

	var payload dto.UpdateCartItemRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "invalid payload given", err)
		return
	}

	cart, err := c.cartService.UpdateCartItem(ctx, userId, uint(id), &payload)
	if err != nil {
		helper.InternalServerError(ctx, "error updating cart", err)
		return
	}

	helper.SuccessResponse(ctx, "item updated successfully", cart)
}

func (c *CartHandler) RemoveCart(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		helper.BadRequestResponse(ctx, "invalid id", err)
		return
	}

	if err := c.cartService.RemoveFromCart(ctx, userId, uint(id)); err != nil {
		helper.InternalServerError(ctx, "error deleting the cart item", err)
		return
	}

	helper.SuccessResponse(ctx, "item removed from cart successfully", nil)
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}
