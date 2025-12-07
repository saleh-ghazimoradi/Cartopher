package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id            uint           `json:"id" gorm:"primaryKey"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Password      string         `json:"-" gorm:"not null"`
	FirstName     string         `json:"first_name" gorm:"not null"`
	LastName      string         `json:"last_name" gorm:"not null"`
	Phone         string         `json:"phone"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	Role          UserRole       `json:"role" gorm:"default:customer"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	RefreshTokens []RefreshToken `json:"-"`
	Orders        []Order        `json:"-"`
	Cart          Cart           `json:"-"`
}

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

type RefreshToken struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	UserId    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"-"`
}
