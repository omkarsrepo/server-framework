package sfk

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type RateLimiterMiddleware struct {
	limiter ratelimit.Limiter
}

func (rl *RateLimiterMiddleware) applyFilter() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		rl.limiter.Take()

		ginCtx.Next()
	}
}

func ApplyRateLimiter() gin.HandlerFunc {
	config := ConfigServiceInstance()
	limiter := ratelimit.New(config.GetInt("rateLimitCallsPerSec"))
	instance := &RateLimiterMiddleware{limiter}

	return instance.applyFilter()
}
