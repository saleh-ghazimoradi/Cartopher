package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func (a *AuthHandler) Register(ctx *gin.Context) {
	var payload dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	registerResponse, err := a.authService.Register(ctx, &payload)
	if err != nil {
		helper.InternalServerError(ctx, "Failed to register", err)
		return
	}

	helper.CreatedResponse(ctx, "user registered successfully", registerResponse)
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var payload dto.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	loginResponse, err := a.authService.Login(ctx, &payload)
	if err != nil {
		helper.InternalServerError(ctx, "Failed to login", err)
		return
	}

	helper.SuccessResponse(ctx, "user logged in successfully", loginResponse)
}

func (a *AuthHandler) RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	refreshTokenResponse, err := a.authService.RefreshToken(ctx, &payload)
	if err != nil {
		helper.InternalServerError(ctx, "Failed to refresh", err)
		return
	}

	helper.SuccessResponse(ctx, "user refreshed successfully", refreshTokenResponse)
}

func (a *AuthHandler) Logout(ctx *gin.Context) {
	var payload dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	if err := a.authService.Logout(ctx, payload.RefreshToken); err != nil {
		helper.InternalServerError(ctx, "Failed to logout", err)
		return
	}

	helper.SuccessResponse(ctx, "user logged out successfully", nil)
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
