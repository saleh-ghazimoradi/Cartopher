package resolver

import (
	"context"
	"errors"
)

const (
	adminRole = "admin"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func GetUserIdFromContext(ctx context.Context) (uint, error) {
	userId := ctx.Value("user_id")

	if userId == nil {
		return 0, ErrUnauthorized
	}

	if id, ok := userId.(uint); ok {
		return id, nil
	}

	return 0, ErrUnauthorized
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value("user_role")

	if userRole == nil {
		return "", ErrUnauthorized
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", ErrUnauthorized
}

func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}
