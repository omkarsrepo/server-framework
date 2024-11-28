// Unpublished Work © 2024

package sfk

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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
	RegisterShutdownHook(cleanup func())
	OverrideCorsWithMiddleware(override bool)
	RegisterMiddlewares(middlewares []gin.HandlerFunc)
	Start(routes func(), database func())
}

type serverService struct {
	cmd                *cobra.Command
	logger             *zerolog.Logger
	config             ConfigService
	router             *gin.Engine
	cleanup            func()
	shouldOverrideCors bool
	middlewares        []gin.HandlerFunc
}

func NewServerService(name, description string) ServerService {
	cobraCmd := &cobra.Command{
		Use:   name,
		Short: description,
	}

	commandsService := newCommandsService(cobraCmd)
	commandsService.registerCommands()

	routerInstance := RouterInstance()
	loggerInstance := LoggerServiceInstance()

	return &serverService{
		cmd:                cobraCmd,
		logger:             loggerInstance.ZeroLogger(),
		config:             ConfigServiceInstance(),
		router:             routerInstance.Router(),
		shouldOverrideCors: false,
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

	middlewareService := newMiddlewareService()
	middlewareService.registerMiddlewares(s.shouldOverrideCors, s.middlewares...)

	customValidators := newCustomValidatorsService()
	customValidators.registerCustomValidators()

	if routes != nil {
		routes()
	}

	if database != nil {
		database()
	}
}

func (s *serverService) startServer() {
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

func (s *serverService) RegisterShutdownHook(cleanup func()) {
	s.cleanup = cleanup
}

func (s *serverService) OverrideCorsWithMiddleware(override bool) {
	s.shouldOverrideCors = override
}

func (s *serverService) RegisterMiddlewares(middlewares []gin.HandlerFunc) {
	s.middlewares = middlewares
}

func (s *serverService) Start(routes func(), database func()) {
	s.cmd.Run = func(_ *cobra.Command, args []string) {
		s.initializeServer(routes, database)
		s.startServer()
	}

	if err := s.cmd.Execute(); err != nil {
		panic(err)
	}
}
