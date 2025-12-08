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

func (u *UserHandler) GetProfile(ctx *gin.Context) {
	userId := ctx.GetUint("user_id")

	profile, err := u.userService.GetProfile(ctx, userId)
	if err != nil {
		helper.NotFoundResponse(ctx, "user not found")
		return
	}

	helper.SuccessResponse(ctx, "User profile retrieved successfully", profile)
}

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
