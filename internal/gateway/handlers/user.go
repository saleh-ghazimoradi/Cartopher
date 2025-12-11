package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/dto"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

// GetProfile docs
// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Failure 404 {object} helper.Response "User not found"
// @Router /users/profile [get]
func (u *UserHandler) GetProfile(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	profile, err := u.userService.GetProfile(ctx, userId)
	if err != nil {
		helper.NotFoundResponse(ctx, "user not found")
		return
	}

	helper.SuccessResponse(ctx, "User profile retrieved successfully", profile)
}

// UpdateProfile docs
// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} helper.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Failure 401 {object} helper.Response "Unauthorized"
// @Router /users/profile [put]
func (u *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	var payload dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		helper.BadRequestResponse(ctx, "invalid request data", err)
		return
	}

	updatedUserProfile, err := u.userService.UpdateProfile(ctx, userId, &payload)
	if err != nil {
		helper.InternalServerError(ctx, "could not update user", err)
		return
	}

	helper.SuccessResponse(ctx, "User profile updated successfully", updatedUserProfile)
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
