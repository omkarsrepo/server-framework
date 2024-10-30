package server_framework

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	routerServiceInstance *routerService
	routerServiceOnce     sync.Once
)

type RouterService interface {
	GetRouter() *gin.Engine
}

type routerService struct {
	*gin.Engine
}

func getRouter() *gin.Engine {
	config := ConfigServiceInstance()
	router := gin.Default()

	if config.GetString("env") == "localhost" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return router
}

func RouterInstance() RouterService {
	routerServiceOnce.Do(func() {
		routerServiceInstance = &routerService{
			Engine: getRouter(),
		}
	})

	return routerServiceInstance
}

func (props *routerService) GetRouter() *gin.Engine {
	return props.Engine
}
