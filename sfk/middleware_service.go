package sfk

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MiddlewareService interface {
	registerMiddlewares(overrideCorsMiddleware bool, middlewares ...gin.HandlerFunc)
}

type middlewareService struct {
	router *gin.Engine
}

func newMiddlewareService() MiddlewareService {
	routerInstance := RouterInstance()
	return &middlewareService{
		router: routerInstance.Router(),
	}
}

func applyCors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}

func (m *middlewareService) registerMiddlewares(overrideCorsMiddleware bool, middlewares ...gin.HandlerFunc) {
	m.router.Use(applyRateLimiter())
	m.router.Use(applyRequestTimeout())
	m.router.Use(applyTraceHeader())
	m.router.Use(applyRequestLoggerMiddleware())

	if !overrideCorsMiddleware {
		m.router.Use(applyCors())
	}

	m.router.Use(middlewares...)
}
