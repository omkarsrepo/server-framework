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
	Run(routes func(), middlewares []gin.HandlerFunc, database func())
}

type serverService struct {
	cmd     *cobra.Command
	logger  *zerolog.Logger
	config  ConfigService
	router  *gin.Engine
	cleanup func()
}

func NewServerService(name, description string) ServerService {
	cobraCmd := &cobra.Command{
		Use:   name,
		Short: description,
	}

	commandsService := NewCommandsService(cobraCmd)
	commandsService.RegisterCommands()

	routerInstance := RouterInstance()
	loggerInstance := LoggerServiceInstance()

	return &serverService{
		cmd:    cobraCmd,
		logger: loggerInstance.GetZeroLogger(),
		config: ConfigServiceInstance(),
		router: routerInstance.GetRouter(),
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

func (s *serverService) initializeServer(routes func(), middlewares []gin.HandlerFunc, database func()) {
	s.setMaxMemoryLimit()

	middlewareService := NewMiddlewareService()
	middlewareService.RegisterMiddlewares(middlewares...)

	customValidators := NewCustomValidatorsService()
	customValidators.RegisterCustomValidators()

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

func (s *serverService) Run(routes func(), middlewares []gin.HandlerFunc, database func()) {
	s.cmd.Run = func(_ *cobra.Command, args []string) {
		s.initializeServer(routes, middlewares, database)
		s.startServer()
	}

	if err := s.cmd.Execute(); err != nil {
		panic(err)
	}
}
