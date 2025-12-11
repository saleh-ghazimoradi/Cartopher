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

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User registration data"
// @Success      201 {object} helper.Response{data=dto.AuthResponse}
// @Failure      400 {object} helper.Response "Invalid request data or user already exists"
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "User login credentials"
// @Success      200 {object} helper.Response{data=dto.AuthResponse} "Login successfully"
// @Failure      401 {object} helper.Response "Invalid credentials"
// @Router       /auth/login [post]
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

// RefreshToken docs
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} helper.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 401 {object} helper.Response "Invalid refresh token"
// @Router /auth/refresh [post]
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

// Logout docs
// @Summary User logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} helper.Response "Logout successful"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Router /auth/logout [post]
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
