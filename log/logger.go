package log

import (
	"os"

	be "github.com/pickjunk/bgo/error"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger a custom logger for bgo, base on zerolog
type Logger struct {
	zerolog.Logger
}

// LogError override zerolog.Err to handle SystemError
func (l *Logger) LogError(err error) *zerolog.Event {
	event := l.Err(err)

	if e, ok := err.(*be.SystemError); ok {
		event = event.Dict("sys_err", e.Event)
	}

	return event
}

// New a logger
func New(component string) *Logger {
	l := zlog.With().Caller().Str("component", component).Logger()

	if os.Getenv("ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		l = l.Output(&lumberjack.Logger{
			Filename:   "runtime/log/app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     20, //days
		})
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return &Logger{l}
}
