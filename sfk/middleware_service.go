package sfk

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MiddlewareService interface {
	RegisterMiddlewares(overrideCorsMiddleware bool, middlewares ...gin.HandlerFunc)
}

type middlewareService struct {
	router *gin.Engine
}

func NewMiddlewareService() MiddlewareService {
	routerInstance := RouterInstance()
	return &middlewareService{
		router: routerInstance.GetRouter(),
	}
}

func applyCors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}

func (m *middlewareService) RegisterMiddlewares(overrideCorsMiddleware bool, middlewares ...gin.HandlerFunc) {
	m.router.Use(applyRateLimiter())
	m.router.Use(applyRequestTimeout())
	m.router.Use(applyTraceHeader())
	m.router.Use(applyRequestLoggerMiddleware())

	if !overrideCorsMiddleware {
		m.router.Use(applyCors())
	}

	m.router.Use(middlewares...)
}
