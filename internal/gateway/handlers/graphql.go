package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

type GraphQLHandler struct {
	handler *handler.Server
}

func (g *GraphQLHandler) GraphqlHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		g.handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (g *GraphQLHandler) PlayGround() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/graphql/")
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (g *GraphQLHandler) PlayGroundPublic() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground (Public)", "/graphql/public/")
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (g *GraphQLHandler) PlayGroundPrivate() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground (private)", "/graphql/")
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func NewGraphQLHandler(handler *handler.Server) *GraphQLHandler {
	return &GraphQLHandler{
		handler: handler,
	}
}
