package log

import (
	"os"

	bc "github.com/pickjunk/bgo/config"
	be "github.com/pickjunk/bgo/error"
	"github.com/pickjunk/zerolog"
)

// Logger a custom logger for bgo, base on zerolog
type Logger struct {
	zerolog.Logger
}

// LogAndPanic an error
// Please always use this in your business code
// It will handle bgo SystemError properly
func (l *Logger) LogAndPanic(err error) *zerolog.Event {
	event := l.Panic().Err(err)

	if e, ok := err.(*be.SystemError); ok {
		event = event.Merge(e.Event)
	}

	return event
}

var inner zerolog.Logger
var outer zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	inner = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	inner = inner.With().Str("component", "bgo.log").Logger()

	logPath := bc.Get("log").String()
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			inner.Fatal().Err(err).Send()
		}
		inner.Info().Str("file", logPath).Msg("log redirect")

		outer = zerolog.New(f)
		outer = outer.With().Str("component", "bgo.log").Logger()
		outer = outer.Level(zerolog.InfoLevel)
	} else {
		outer = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
		outer = outer.Level(zerolog.DebugLevel)
	}

	outer = outer.Hook(callerHook{})
}

// New a logger
func New(component string) *Logger {
	l := outer.With().Str("component", component).Logger()
	return &Logger{l}
}

var (
	// Dict creates an Event to be used with the *Event.Dict or *Event.Merge method.
	Dict = zerolog.Dict
)
