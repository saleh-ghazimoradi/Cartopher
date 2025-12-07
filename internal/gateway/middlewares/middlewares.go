package middlewares

import "github.com/gin-gonic/gin"

type Middlewares struct{}

func (m *Middlewares) CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-control-Allow-Headers", "Content-Type, Accept, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

func NewMiddlewares() *Middlewares {
	return &Middlewares{}
}
