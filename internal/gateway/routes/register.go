package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/saleh-ghazimoradi/Cartopher/docs"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Register struct {
	middlewares  *middlewares.Middlewares
	healthRoute  *HealthRoutes
	authRoute    *AuthRoutes
	userRoute    *UserRoutes
	productRoute *ProductRoutes
	cartRoute    *CartRoutes
	orderRoute   *OrderRoutes
	graphqlRoute *GraphQLRoutes
}

type Options func(*Register)

func WithHealthRoute(healthRoute *HealthRoutes) Options {
	return func(r *Register) {
		r.healthRoute = healthRoute
	}
}

func WithAuthRoute(authRoute *AuthRoutes) Options {
	return func(r *Register) {
		r.authRoute = authRoute
	}
}

func WithUserRoute(userRoute *UserRoutes) Options {
	return func(r *Register) {
		r.userRoute = userRoute
	}
}

func WithProductRoute(productRoute *ProductRoutes) Options {
	return func(r *Register) {
		r.productRoute = productRoute
	}
}

func WithCartRoute(cartRoute *CartRoutes) Options {
	return func(r *Register) {
		r.cartRoute = cartRoute
	}
}

func WithOrderRoute(orderRoute *OrderRoutes) Options {
	return func(r *Register) {
		r.orderRoute = orderRoute
	}
}

func WithMiddlewares(middlewares *middlewares.Middlewares) Options {
	return func(r *Register) {
		r.middlewares = middlewares
	}
}

func WithGraphqlRoute(graphqlRoute *GraphQLRoutes) Options {
	return func(r *Register) {
		r.graphqlRoute = graphqlRoute
	}
}

func (r *Register) RegisterRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(r.middlewares.CorsMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Static("/docs", "./docs")

	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	r.healthRoute.HealthRoute(router)
	r.authRoute.AuthRoute(router)
	r.userRoute.UserRoute(router)
	r.productRoute.ProductRoute(router)
	r.cartRoute.cartRoute(router)
	r.orderRoute.OrderRoute(router)
	r.graphqlRoute.GraphQLRoute(router)
	return router
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, f := range opts {
		f(r)
	}
	return r
}
