// Unpublished Work © 2024

package types

import "github.com/gin-gonic/gin"

type ServerOptions struct {
	Routes                         func()
	Database                       func()
	ShutdownHook                   func()
	ShouldOverrideCORSMiddleware   bool
	ShouldDisableGzipCompression   bool
	ExcludePathsForGzipCompression []string
	Middlewares                    []gin.HandlerFunc
	SkipRateLimiterMiddleware      bool
	SkipRequestTimeoutMiddleware   bool
	SkipTraceHeaderMiddleware      bool
	SkipRequestLoggerMiddleware    bool
	DisablePprof                   bool
}
