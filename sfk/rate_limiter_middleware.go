package sfk

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type rateLimiterMiddleware struct {
	limiter ratelimit.Limiter
}

func (rl *rateLimiterMiddleware) applyFilter() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		rl.limiter.Take()

		ginCtx.Next()
	}
}

func applyRateLimiter() gin.HandlerFunc {
	config := ConfigServiceInstance()
	limiter := ratelimit.New(config.GetInt("rateLimitCallsPerSec"))
	instance := &rateLimiterMiddleware{limiter}

	return instance.applyFilter()
}
