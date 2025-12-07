package dto

import "time"

type AddToCartRequest struct {
	ProductId uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	Id        uint               `json:"id"`
	UserId    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64
}

type CartItemResponse struct {
	Id       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
}

type OrderResponse struct {
	Id          uint                `json:"id"`
	UserId      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   time.Time           `json:"created_at"`
}

type OrderItemResponse struct {
	Id       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}
