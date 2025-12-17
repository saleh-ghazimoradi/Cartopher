package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"github.com/saleh-ghazimoradi/Cartopher/utils"
)

type Authentication struct {
	config *config.Config
}

func (a *Authentication) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			helper.UnauthorizedResponse(ctx, "Authorization header is required")
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helper.UnauthorizedResponse(ctx, "invalid authorization header format")
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], a.config.JWT.Secret)
		if err != nil {
			helper.UnauthorizedResponse(ctx, "invalid token")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserId)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

func (a *Authentication) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			helper.ForbiddenResponse(ctx, "You are not authorized to access this resource")
			ctx.Abort()
			return
		}

		if role != string(domain.UserRoleAdmin) {
			helper.ForbiddenResponse(ctx, "You are not authorized to access this resource")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (a *Authentication) GraphqlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")
		userRole, _ := c.Get("user_role")

		ctx := context.WithValue(c.Request.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "user_email", userEmail)
		ctx = context.WithValue(ctx, "user_role", userRole)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func NewAuthentication(config *config.Config) *Authentication {
	return &Authentication{
		config: config,
	}
}
