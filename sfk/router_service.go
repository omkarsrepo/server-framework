// Unpublished Work © 2024

package sfk

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

func registerHealthPingEndpoint(router *gin.Engine) {
	router.GET("/health/IhEaf/ping", func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusNoContent, nil)
	})
}

func getRouter() *gin.Engine {
	config := ConfigServiceInstance()
	router := gin.Default()

	if config.GetString("env") != "prod" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	registerHealthPingEndpoint(router)

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
