package sfk

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareService interface {
	RegisterMiddlewares(middlewares ...gin.HandlerFunc)
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

func (m *middlewareService) RegisterMiddlewares(middlewares ...gin.HandlerFunc) {
	m.router.Use(ApplyRateLimiter())
	m.router.Use(ApplyRequestTimeout())
	m.router.Use(ApplyTraceHeader())
	m.router.Use(middlewares...)
}
