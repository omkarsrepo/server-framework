package sf

import (
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
	Info() *zerolog.Event
	Error() *zerolog.Event
	Err(err error) *zerolog.Event
	Fatal() *zerolog.Event
	Panic() *zerolog.Event
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

func (props *loggerService) Info() *zerolog.Event {
	return props.Logger.Info()
}

func (props *loggerService) Error() *zerolog.Event {
	return props.Logger.Error()
}

func (props *loggerService) Err(err error) *zerolog.Event {
	return props.Logger.Err(err)
}

func (props *loggerService) Fatal() *zerolog.Event {
	return props.Logger.Fatal()
}

func (props *loggerService) Panic() *zerolog.Event {
	return props.Logger.Panic()
}

func (props *loggerService) GetZeroLogger() *zerolog.Logger {
	return props.Logger
}
