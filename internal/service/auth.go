package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
	"github.com/saleh-ghazimoradi/Cartopher/utils"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	cfg            *config.Config
	userRepository repository.UserRepository
	cartRepository repository.CartRepository
}

func (a *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := a.userRepository.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("you cannot register with this user")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      domain.UserRoleCustomer,
	}

	if err := a.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	cart := &domain.Cart{UserId: user.Id}
	if err := a.cartRepository.CreateCart(ctx, cart); err != nil {
		return nil, err
	}

	return a.generateAuthResponse(ctx, user)

}

func (a *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := a.userRepository.GetUserByEmailAndActive(ctx, req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return a.generateAuthResponse(ctx, user)
}

func (a *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claim, err := utils.ValidateToken(req.RefreshToken, a.cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := a.userRepository.GetValidRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	user, err := a.userRepository.GetUserById(ctx, claim.UserId)
	if err != nil {
		return nil, err
	}

	if err := a.userRepository.DeleteRefreshTokenById(ctx, refreshToken.Id); err != nil {
		return nil, err
	}

	return a.generateAuthResponse(ctx, user)
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	return a.userRepository.DeleteRefreshToken(ctx, refreshToken)
}

func (a *authService) generateAuthResponse(ctx context.Context, user *domain.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateToken(
		a.cfg,
		user.Id,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenDomain := &domain.RefreshToken{
		UserId:    user.Id,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.cfg.JWT.RefreshTokenExpires),
	}

	if err := a.userRepository.CreateRefreshToken(ctx, refreshTokenDomain); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			Id:        user.Id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuthService(cfg *config.Config, userRepository repository.UserRepository, cartRepository repository.CartRepository) AuthService {
	return &authService{
		cfg:            cfg,
		userRepository: userRepository,
		cartRepository: cartRepository,
	}
}
