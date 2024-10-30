package sf

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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
	logger  LoggerService
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

	return &serverService{
		cmd:    cobraCmd,
		logger: LoggerServiceInstance(),
		config: ConfigServiceInstance(),
		router: routerInstance.GetRouter(),
	}
}

func (props *serverService) cleanupOnShutdown() {
	go func() {
		Cache().Close()

		if props.cleanup != nil {
			props.cleanup()
		}
	}()
}

func (props *serverService) setMaxMemoryLimit() {
	if props.config.GetString("env") != "localhost" {
		maxMemoryLimit := props.config.GetInt64("maxMemoryLimitInMB")

		debug.SetMemoryLimit(maxMemoryLimit * 1 << 20)
	}
}

func (props *serverService) shutdownGracefully(server *http.Server) {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	props.logger.Info().Msg("Received shutdown server event...")

	gracefulShutdownSecs := props.config.GetInt("gracefulShutdownSecs")
	gracefulShutdown := time.Duration(gracefulShutdownSecs) * time.Second

	props.logger.Info().Msgf("Server Shutdown timeout of %s seconds...", gracefulShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdown)
	defer cancel()

	props.cleanupOnShutdown()

	if err := server.Shutdown(ctx); err != nil {
		props.logger.Fatal().Msgf("Server failed to gracefully shutdown before timeout: %+v", err)
	}

	<-ctx.Done()
	props.logger.Info().Msgf("Server Shutdown timeout of %s seconds completed successfully. Server Exited!", gracefulShutdown)
}

func (props *serverService) initializeServer(routes func(), middlewares []gin.HandlerFunc, database func()) {
	props.setMaxMemoryLimit()

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

func (props *serverService) startServer() {
	port := props.config.GetString("port")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: props.router.Handler(),
	}

	go func() {
		props.logger.Info().Msgf("Server running successfully on port %s", port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			props.logger.Fatal().Msgf("Failed to start server... listen: %+v\n", err)
		}
	}()

	props.shutdownGracefully(server)
}

func (props *serverService) RegisterShutdownHook(cleanup func()) {
	props.cleanup = cleanup
}

func (props *serverService) Run(routes func(), middlewares []gin.HandlerFunc, database func()) {
	props.cmd.Run = func(_ *cobra.Command, args []string) {
		props.initializeServer(routes, middlewares, database)
		props.startServer()
	}

	if err := props.cmd.Execute(); err != nil {
		panic(err)
	}
}
