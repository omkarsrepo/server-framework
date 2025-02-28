// Unpublished Work © 2024

package sfk

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type MiddlewareService interface {
	registerMiddlewares(middlewares ...gin.HandlerFunc)
}

type middlewareOptions struct {
	overrideCorsMiddleware         bool
	disableGzipCompression         bool
	excludePathsForGzipCompression []string
	skipRateLimiterMiddleware      bool
	skipRequestTimeoutMiddleware   bool
	skipTraceHeaderMiddleware      bool
	skipRequestLoggerMiddleware    bool
}

type middlewareService struct {
	router  *gin.Engine
	options *middlewareOptions
}

func newMiddlewareService(options *middlewareOptions) MiddlewareService {
	routerInstance := RouterInstance()
	return &middlewareService{
		router:  routerInstance.Router(),
		options: options,
	}
}

func applyCors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}

func (m *middlewareService) registerMiddlewares(middlewares ...gin.HandlerFunc) {
	if !m.options.skipRateLimiterMiddleware {
		m.router.Use(applyRateLimiter())
	}

	if !m.options.skipRequestTimeoutMiddleware {
		m.router.Use(ApplyRequestTimeout())
	}

	if !m.options.skipTraceHeaderMiddleware {
		m.router.Use(applyTraceHeader())
	}

	if !m.options.skipRequestLoggerMiddleware {
		m.router.Use(applyRequestLoggerMiddleware())
	}

	if !m.options.overrideCorsMiddleware {
		m.router.Use(applyCors())
	}

	if !m.options.disableGzipCompression {
		m.router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions(m.options.excludePathsForGzipCompression)))
	}

	m.router.Use(applyStringReqBody())

	m.router.Use(middlewares...)
}
