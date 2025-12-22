package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/internal/helper"
	"net/http"
	"strconv"
	"time"
)

type Middlewares struct {
	cfg         *config.Config
	rateLimiter helper.RateLimiter
}

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

func (m *Middlewares) RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var key string

		if userID, exists := ctx.Get("user_id"); exists {
			key = "rate_limit:user:" + userID.(string)
		} else {
			key = "rate_limit:ip:" + ctx.ClientIP()
		}

		remaining, retryAfter, allowed, err := m.rateLimiter.Allow(ctx.Request.Context(), key, m.cfg.Redis.RPM)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Header("RateLimit_Remaining", strconv.Itoa(remaining))

		if !allowed {
			ctx.Header(
				"RateLimit-RetryAfter",
				strconv.Itoa(int(retryAfter/time.Second)),
			)

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		if err != nil {
			ctx.Next()
			return
		}

	}
}

func NewMiddlewares(cfg *config.Config, rateLimiter helper.RateLimiter) *Middlewares {
	return &Middlewares{
		cfg:         cfg,
		rateLimiter: rateLimiter,
	}
}
