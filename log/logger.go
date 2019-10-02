package log

import (
	"os"

	bc "github.com/pickjunk/bgo/config"
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
	l := zlog.With().Str("component", component).Logger()

	logPath := bc.Get("log.path").String()
	if logPath != "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		maxSize := bc.Get("log.maxSize").Uint()
		if maxSize == 0 {
			maxSize = 10
		}
		maxBackups := bc.Get("log.maxBackups").Uint()
		if maxBackups == 0 {
			maxBackups = 10
		}
		maxAge := bc.Get("log.maxAge").Uint()
		if maxSize == 0 {
			maxAge = 20
		}
		l = l.Output(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    int(maxSize), // megabytes
			MaxBackups: int(maxBackups),
			MaxAge:     int(maxAge), //days
		})
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	l = l.Hook(callerHook{})

	return &Logger{l}
}
