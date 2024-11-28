// Unpublished Work © 2024

package sfk

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	routerServiceInstance *routerService
	routerServiceOnce     sync.Once
)

type RouterService interface {
	Router() *gin.Engine
}

type routerService struct {
	*gin.Engine
}

func getRouter() *gin.Engine {
	config := ConfigServiceInstance()
	router := gin.Default()

	if config.GetString("env") != "prod" {
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

func (r *routerService) Router() *gin.Engine {
	return r.Engine
}
