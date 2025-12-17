package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/graph"
	"github.com/saleh-ghazimoradi/Cartopher/graph/resolver"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQL struct {
	authService    service.AuthService
	userService    service.UserService
	productService service.ProductService
	cartService    service.CartService
	orderService   service.OrderService
}

type Option func(*GraphQL)

func WithAuthService(authService service.AuthService) Option {
	return func(g *GraphQL) {
		g.authService = authService
	}
}

func WithUserService(userService service.UserService) Option {
	return func(g *GraphQL) {
		g.userService = userService
	}
}

func WithProductService(productService service.ProductService) Option {
	return func(g *GraphQL) {
		g.productService = productService
	}
}

func WithCartService(cartService service.CartService) Option {
	return func(g *GraphQL) {
		g.cartService = cartService
	}
}

func WithOrderService(orderService service.OrderService) Option {
	return func(g *GraphQL) {
		g.orderService = orderService
	}
}

func (g *GraphQL) GraphqlHandler() *handler.Server {
	rvr := resolver.NewResolver(
		resolver.WithAuthService(g.authService),
		resolver.WithUserService(g.userService),
		resolver.WithProductService(g.productService),
		resolver.WithCartService(g.cartService),
		resolver.WithOrderService(g.orderService),
	)

	schema := graph.NewExecutableSchema(graph.Config{Resolvers: rvr})

	srv := handler.New(schema)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func (g *GraphQL) GraphqlHandlerWrapper() gin.HandlerFunc {
	h := g.GraphqlHandler()
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (g *GraphQL) PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/graphql/")
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func NewGraphQL(opts ...Option) *GraphQL {
	g := &GraphQL{}
	for _, o := range opts {
		o(g)
	}
	return g
}
