package sfk

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

func extractTraceId(ginCtx *gin.Context) (string, bool) {
	traceId := ""

	if ginCtx != nil {
		traceId = ginCtx.GetString("TRACE_ID")
		if traceId != "" {
			return traceId, true
		}
	}

	return traceId, false
}

func (l *loggerService) Info(ginCtx *gin.Context) *zerolog.Event {
	traceId, ok := extractTraceId(ginCtx)
	if ok {
		logger := l.Logger.With().Str("traceId", traceId).Logger()

		return logger.Info()
	}

	return l.Logger.Info()
}

func (l *loggerService) Error(ginCtx *gin.Context) *zerolog.Event {
	traceId, ok := extractTraceId(ginCtx)
	if ok {
		logger := l.Logger.With().Str("traceId", traceId).Logger()

		return logger.Error()
	}

	return l.Logger.Error()
}

func (l *loggerService) Err(ginCtx *gin.Context, err error) *zerolog.Event {
	traceId, ok := extractTraceId(ginCtx)
	if ok {
		logger := l.Logger.With().Str("traceId", traceId).Logger()

		return logger.Err(err)
	}

	return l.Logger.Err(err)
}

func (l *loggerService) Fatal(ginCtx *gin.Context) *zerolog.Event {
	traceId, ok := extractTraceId(ginCtx)
	if ok {
		logger := l.Logger.With().Str("traceId", traceId).Logger()

		return logger.Fatal()
	}

	return l.Logger.Fatal()
}

func (l *loggerService) Panic(ginCtx *gin.Context) *zerolog.Event {
	traceId, ok := extractTraceId(ginCtx)
	if ok {
		logger := l.Logger.With().Str("traceId", traceId).Logger()

		return logger.Panic()
	}

	return l.Logger.Panic()
}

func (l *loggerService) GetZeroLogger() *zerolog.Logger {
	return l.Logger
}
