package server_framework

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type RateLimiterMiddleware struct {
	limiter ratelimit.Limiter
}

func (props *RateLimiterMiddleware) applyFilter() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		props.limiter.Take()

		ginCtx.Next()
	}
}

func ApplyRateLimiter() gin.HandlerFunc {
	config := ConfigServiceInstance()
	limiter := ratelimit.New(config.GetInt("rateLimitCallsPerSec"))
	instance := &RateLimiterMiddleware{limiter}

	return instance.applyFilter()
}
