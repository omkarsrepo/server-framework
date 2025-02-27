// Unpublished Work © 2024

package sfk

import (
	"context"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sfk/boom"
	"github.com/omkarsrepo/server-framework/sfk/json"
	"github.com/omkarsrepo/server-framework/sfk/types"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

type ServerService interface {
	Start()
}

type serverService struct {
	cmd                            *cobra.Command
	logger                         *zerolog.Logger
	config                         ConfigService
	router                         *gin.Engine
	cleanup                        func()
	shouldOverrideCors             bool
	middlewares                    []gin.HandlerFunc
	routes                         func()
	database                       func()
	disableGzipCompression         bool
	excludePathsForGzipCompression []string
	skipRateLimiterMiddleware      bool
	skipRequestTimeoutMiddleware   bool
	skipTraceHeaderMiddleware      bool
	skipRequestLoggerMiddleware    bool
	disablePprof                   bool
}

func NewServerService(name, description string, options *types.ServerOptions) ServerService {
	cobraCmd := &cobra.Command{
		Use:   name,
		Short: description,
	}

	commandsService := newCommandsService(cobraCmd)
	commandsService.registerCommands()

	routerInstance := RouterInstance()
	loggerInstance := LoggerServiceInstance()

	return &serverService{
		cmd:                            cobraCmd,
		logger:                         loggerInstance.ZeroLogger(),
		config:                         ConfigServiceInstance(),
		router:                         routerInstance.Router(),
		shouldOverrideCors:             options.ShouldOverrideCORSMiddleware,
		middlewares:                    options.Middlewares,
		cleanup:                        options.ShutdownHook,
		routes:                         options.Routes,
		database:                       options.Database,
		disableGzipCompression:         options.ShouldDisableGzipCompression,
		excludePathsForGzipCompression: options.ExcludePathsForGzipCompression,
		skipRateLimiterMiddleware:      options.SkipRateLimiterMiddleware,
		skipRequestTimeoutMiddleware:   options.SkipRequestTimeoutMiddleware,
		skipTraceHeaderMiddleware:      options.SkipTraceHeaderMiddleware,
		skipRequestLoggerMiddleware:    options.SkipRequestLoggerMiddleware,
		disablePprof:                   options.DisablePprof,
	}
}

func (s *serverService) cleanupOnShutdown() {
	go func() {
		Cache().Close()

		if s.cleanup != nil {
			s.cleanup()
		}
	}()
}

func (s *serverService) setMaxMemoryLimit() {
	if s.config.GetString("env") != "localhost" {
		maxMemoryLimit := s.config.GetInt64("maxMemoryLimitInMB")

		debug.SetMemoryLimit(maxMemoryLimit * 1 << 20)
	}
}

func (s *serverService) shutdownGracefully(server *http.Server) {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	s.logger.Info().Msg("Received shutdown server event...")

	gracefulShutdownSecs := s.config.GetInt("gracefulShutdownSecs")
	gracefulShutdown := time.Duration(gracefulShutdownSecs) * time.Second

	s.logger.Info().Msgf("Server Shutdown timeout of %s seconds...", gracefulShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdown)
	defer cancel()

	s.cleanupOnShutdown()

	if err := server.Shutdown(ctx); err != nil {
		s.logger.Fatal().Msgf("Server failed to gracefully shutdown before timeout: %+v", err)
	}

	<-ctx.Done()
	s.logger.Info().Msgf("Server Shutdown timeout of %s seconds completed successfully. Server Exited!", gracefulShutdown)
}

func (s *serverService) initializeServer(routes func(), database func()) {
	s.setMaxMemoryLimit()

	middlewareService := newMiddlewareService(&middlewareOptions{
		overrideCorsMiddleware:         s.shouldOverrideCors,
		disableGzipCompression:         s.disableGzipCompression,
		excludePathsForGzipCompression: s.excludePathsForGzipCompression,
		skipRateLimiterMiddleware:      s.skipRateLimiterMiddleware,
		skipRequestTimeoutMiddleware:   s.skipRequestTimeoutMiddleware,
		skipTraceHeaderMiddleware:      s.skipTraceHeaderMiddleware,
		skipRequestLoggerMiddleware:    s.skipRequestLoggerMiddleware,
	})

	middlewareService.registerMiddlewares(s.middlewares...)

	customValidators := newCustomValidatorsService()
	customValidators.registerCustomValidators()

	if routes != nil {
		routes()
	}

	if database != nil {
		database()
	}
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

func (s *serverService) startServer() {
	if !s.disablePprof {
		enablePprof(s.router)
	}

	port := s.config.GetString("port")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router.Handler(),
	}

	go func() {
		s.logger.Info().Msgf("Server running successfully on port %s", port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal().Msgf("Failed to start server... listen: %+v\n", err)
		}
	}()

	s.shutdownGracefully(server)
}

func (s *serverService) Start() {
	s.cmd.Run = func(_ *cobra.Command, args []string) {
		s.initializeServer(s.routes, s.database)
		s.startServer()
	}

	if err := s.cmd.Execute(); err != nil {
		panic(err)
	}
}
