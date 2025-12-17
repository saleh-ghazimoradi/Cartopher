package resolver

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Cartopher/utils"
)

const (
	adminRole = "admin"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func GetUserIdFromContext(ctx context.Context) (uint, error) {
	userId := ctx.Value(utils.UserIdKey)

	if userId == nil {
		return 0, ErrUnauthorized
	}

	if id, ok := userId.(uint); ok {
		return id, nil
	}

	return 0, ErrUnauthorized
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value(utils.UserRoleKey)

	if userRole == nil {
		return "", ErrUnauthorized
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", ErrUnauthorized
}

func getPagingNumbers(page, limit *int) (int, int) {
	p, l := 0, 0

	if page != nil {
		p = *page
	}

	if limit != nil {
		l = *limit
	}

	if p <= 0 {
		p = 1
	}

	if l <= 0 {
		l = 10
	}

	return p, l
}

func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}
