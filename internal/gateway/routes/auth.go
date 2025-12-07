package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
)

type AuthRoutes struct {
	authHandler *handlers.AuthHandler
}

func (a *AuthRoutes) AuthRoute(router *gin.Engine) {
	v1 := router.Group("/v1")
	auth := v1.Group("/auth")
	auth.POST("/register", a.authHandler.Register)
	auth.POST("/login", a.authHandler.Login)
	auth.POST("/refresh", a.authHandler.RefreshToken)
	auth.POST("/logout", a.authHandler.Logout)
}

func NewAuthRoutes(authHandler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{authHandler: authHandler}
}
