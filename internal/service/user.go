package service

import (
	"context"

	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userId uint) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userId uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) GetProfile(ctx context.Context, userId uint) (*dto.UserResponse, error) {
	user, err := u.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
	}, nil
}

func (u *userService) UpdateProfile(ctx context.Context, userId uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := u.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	if err := u.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return u.GetProfile(ctx, userId)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
