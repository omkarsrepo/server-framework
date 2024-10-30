package server_framework

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

func (props *middlewareService) RegisterMiddlewares(middlewares ...gin.HandlerFunc) {
	props.router.Use(ApplyRateLimiter())
	props.router.Use(ApplyRequestTimeout())
	props.router.Use(middlewares...)
}
