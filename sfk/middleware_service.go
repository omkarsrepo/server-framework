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
	m.router.Use(applyRateLimiter())
	m.router.Use(applyRequestTimeout())
	m.router.Use(applyTraceHeader())
	m.router.Use(applyRequestLoggerMiddleware())

	if !m.options.overrideCorsMiddleware {
		m.router.Use(applyCors())
	}

	if !m.options.disableGzipCompression {
		m.router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions(m.options.excludePathsForGzipCompression)))
	}

	m.router.Use(middlewares...)
}
