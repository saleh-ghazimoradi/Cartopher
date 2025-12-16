package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
)

type CartHandler struct {
	cartService service.CartService
}

// GetCart docs
// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 404 {object} helper.Response "Cart not found"
// @Router /cart [get]
func (c *CartHandler) GetCart(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	cart, err := c.cartService.GetCart(ctx, userId)
	if err != nil {
		helper.NotFoundResponse(ctx, "Cart not found")
		return
	}

	helper.SuccessResponse(ctx, "Cart retrieved successfully", cart)
}

// AddToCart docs
// @Summary Add item to cart
// @Description Add a product to the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Item to add to cart"
// @Success 200 {object} helper.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 400 {object} helper.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Router /cart/items [post]
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

// UpdateCart docs
// @Summary Update cart item quantity
// @Description Update the quantity of an item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} helper.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} helper.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Router /cart/items/{id} [put]
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

// RemoveCart docs
// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} helper.Response "Item removed from cart successfully"
// @Failure 400 {object} helper.Response "Invalid cart item ID"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Router /cart/items/{id} [delete]
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
