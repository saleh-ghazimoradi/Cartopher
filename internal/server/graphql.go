package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/saleh-ghazimoradi/Cartopher/graph"
	"github.com/saleh-ghazimoradi/Cartopher/graph/resolver"
	"github.com/vektah/gqlparser/v2/ast"
)

type Graphql struct {
	resolver *resolver.Resolver
}

func (g *Graphql) Connect() *handler.Server {
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: g.resolver})
	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv

}

func NewGraphql(resolver *resolver.Resolver) *Graphql {
	return &Graphql{
		resolver: resolver,
	}
}
