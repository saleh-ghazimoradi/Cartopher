package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type UserRoutes struct {
	userHandler    *handlers.UserHandler
	authMiddleware *middlewares.Authentication
}

func (u *UserRoutes) UserRoute(router *gin.Engine) {
	v1 := router.Group("/v1")
	protected := v1.Group("/")
	protected.Use(u.authMiddleware.Authenticate())
	user := protected.Group("/users")
	user.GET("/profile", u.userHandler.GetProfile)
	user.PUT("/profile", u.userHandler.UpdateProfile)
}

func NewUserRoutes(userHandler *handlers.UserHandler, authMiddleware *middlewares.Authentication) *UserRoutes {
	return &UserRoutes{
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}
