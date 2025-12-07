package domain

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Products []Product `json:"-"`
}

type Product struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	CategoryId  uint           `json:"category_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Category   Category       `json:"Category"`
	Images     []ProductImage `json:"images"`
	OrderItems []OrderItem    `json:"-"`
	CartItems  []CartItem     `json:"-"`
}

type ProductImage struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	ProductId uint           `json:"product_id" gorm:"not null"`
	URL       string         `json:"url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	IsPrimary bool           `json:"is_primary" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Product Product `json:"-"`
}
