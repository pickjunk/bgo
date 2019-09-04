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

// Event a custom event for bgo, base on zerolog
type Event struct {
	*zerolog.Event
}

// Panic custom panic to return a custom event
func (l *Logger) Panic() *Event {
	return &Event{l.Logger.Panic()}
}

// Err custom Logger.Err to handle SystemError
func (l *Logger) Err(err error) *zerolog.Event {
	event := l.Logger.Err(err)

	if e, ok := err.(*be.SystemError); ok {
		event = event.Dict("inner", e.Event)
	}

	return event
}

// Err custom Event.Err to handle SystemError
func (l *Event) Err(err error) *zerolog.Event {
	event := l.Event.Err(err)

	if e, ok := err.(*be.SystemError); ok {
		event = event.Dict("inner", e.Event)
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
