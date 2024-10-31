package sf

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var (
	loggerServiceInstance *loggerService
	loggerServiceOnce     sync.Once
)

type LoggerService interface {
	GetZeroLogger() *zerolog.Logger
	Info(ginCtx *gin.Context) *zerolog.Event
	Error(ginCtx *gin.Context) *zerolog.Event
	Err(ginCtx *gin.Context, err error) *zerolog.Event
	Fatal(ginCtx *gin.Context) *zerolog.Event
	Panic(ginCtx *gin.Context) *zerolog.Event
}

type loggerService struct {
	*zerolog.Logger
}

func getLogger() *zerolog.Logger {
	environment := ConfigServiceInstance().GetString("env")

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Str("env", environment).
		Caller().
		Timestamp().Logger()

	return &logger
}

func LoggerServiceInstance() LoggerService {
	loggerServiceOnce.Do(func() {
		loggerServiceInstance = &loggerService{
			Logger: getLogger(),
		}
	})

	return loggerServiceInstance
}

func (props *loggerService) Info(ginCtx *gin.Context) *zerolog.Event {
	traceId := ginCtx.GetString("TRACE_ID")
	logger := props.Logger.With().Str("traceId", traceId).Logger()

	return logger.Info()
}

func (props *loggerService) Error(ginCtx *gin.Context) *zerolog.Event {
	traceId := ginCtx.GetString("TRACE_ID")
	logger := props.Logger.With().Str("traceId", traceId).Logger()

	return logger.Error()
}

func (props *loggerService) Err(ginCtx *gin.Context, err error) *zerolog.Event {
	traceId := ginCtx.GetString("TRACE_ID")
	logger := props.Logger.With().Str("traceId", traceId).Logger()

	return logger.Err(err)
}

func (props *loggerService) Fatal(ginCtx *gin.Context) *zerolog.Event {
	traceId := ginCtx.GetString("TRACE_ID")
	logger := props.Logger.With().Str("traceId", traceId).Logger()

	return logger.Fatal()
}

func (props *loggerService) Panic(ginCtx *gin.Context) *zerolog.Event {
	traceId := ginCtx.GetString("TRACE_ID")
	logger := props.Logger.With().Str("traceId", traceId).Logger()

	return logger.Panic()
}

func (props *loggerService) GetZeroLogger() *zerolog.Logger {
	return props.Logger
}
