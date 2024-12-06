// Unpublished Work © 2024

package sfk

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sfk/boom"
	"github.com/omkarsrepo/server-framework/sfk/json"
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

func enablePprof(router *gin.Engine) {
	secretService := SecretServiceInstance()

	pprofEndpoint := router.Group("/metrics", func(ginCtx *gin.Context) {
		authToken, exp := json.ExtractAuthorization(ginCtx)
		if exp != nil {
			Abort(ginCtx, exp)
			return
		}

		val, exp := secretService.ValueOf("pprofSecret")
		if exp != nil {
			Abort(ginCtx, exp)
			return
		}

		if authToken != val {
			Abort(ginCtx, boom.Unauthorized("Invalid authToken for authorization header"))
			return
		}

		ginCtx.Next()
	})

	pprof.RouteRegister(pprofEndpoint, "pprof")
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
	enablePprof(router)

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
