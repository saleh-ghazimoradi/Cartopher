package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
)

type Register struct {
	middlewares  *middlewares.Middlewares
	healthRoute  *HealthRoutes
	authRoute    *AuthRoutes
	userRoute    *UserRoutes
	productRoute *ProductRoutes
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

func WithMiddlewares(middlewares *middlewares.Middlewares) Options {
	return func(r *Register) {
		r.middlewares = middlewares
	}
}

func (r *Register) RegisterRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(r.middlewares.CorsMiddleware())

	r.healthRoute.HealthRoute(router)
	r.authRoute.AuthRoute(router)
	r.userRoute.UserRoute(router)
	r.productRoute.ProductRoute(router)
	return router
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, f := range opts {
		f(r)
	}
	return r
}
