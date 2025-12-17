package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type GraphQLRoutes struct {
	graphqlHandler *handlers.GraphQLHandler
	authMiddleware *middlewares.Authentication
}

func (g *GraphQLRoutes) GraphQLRoute(router *gin.Engine) {
	router.GET("/playground", g.graphqlHandler.PlayGround())
	router.GET("/playground/public", g.graphqlHandler.PlayGroundPublic())
	router.GET("/playground/protected", g.graphqlHandler.PlayGroundPrivate())

	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(g.authMiddleware.GraphqlMiddleware())
	graphqlPublic.POST("/", g.graphqlHandler.GraphqlHandler())

	graphqlProtected := router.Group("/graphql")
	graphqlProtected.Use(g.authMiddleware.Authenticate())
	graphqlProtected.Use(g.authMiddleware.GraphqlMiddleware())
	graphqlProtected.POST("/", g.graphqlHandler.GraphqlHandler())

}

func NewGraphQLRoutes(handler *handlers.GraphQLHandler, authMiddleware *middlewares.Authentication) *GraphQLRoutes {
	return &GraphQLRoutes{
		graphqlHandler: handler,
		authMiddleware: authMiddleware,
	}
}
