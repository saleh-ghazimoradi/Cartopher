package helper

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, rpm int) (remaining int, retryAfter time.Duration, allowed bool, err error)
}

type rateLimiter struct {
	limiter *redis_rate.Limiter
}

func (r *rateLimiter) Allow(ctx context.Context, key string, rpm int) (remaining int, retryAfter time.Duration, allowed bool, err error) {
	res, err := r.limiter.Allow(ctx, key, redis_rate.PerMinute(rpm))
	if err != nil {
		return 0, 0, false, err
	}

	return res.Remaining, res.RetryAfter, res.Allowed > 0, nil
}

func NewRateLimiter(limiter *redis_rate.Limiter) RateLimiter {
	return &rateLimiter{
		limiter: limiter,
	}
}
