package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserById(ctx context.Context, id uint) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByEmailAndActive(ctx context.Context, email string, isActive bool) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error

	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetValidRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenById(ctx context.Context, id uint) error
}

type userRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return u.dbWrite.WithContext(ctx).Create(user).Error
}

func (u *userRepository) GetUserById(ctx context.Context, id uint) (*domain.User, error) {
	var user *domain.User
	if err := u.dbRead.WithContext(ctx).First(&user, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user *domain.User
	if err := u.dbRead.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) GetUserByEmailAndActive(ctx context.Context, email string, isActive bool) (*domain.User, error) {
	var user *domain.User
	if err := u.dbRead.WithContext(ctx).Where("email = ? AND is_active = ?", email, isActive).First(&user).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.dbWrite.WithContext(ctx).Save(user).Error
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	return u.dbWrite.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (u *userRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return u.dbWrite.WithContext(ctx).Create(token).Error
}

func (u *userRepository) GetValidRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken *domain.RefreshToken
	if err := u.dbRead.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (u *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return u.dbWrite.WithContext(ctx).Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (u *userRepository) DeleteRefreshTokenById(ctx context.Context, id uint) error {
	return u.dbWrite.WithContext(ctx).Delete(&domain.RefreshToken{}, id).Error
}

func NewAuthRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
