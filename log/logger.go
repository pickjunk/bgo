package log

import (
	"os"

	bc "github.com/pickjunk/bgo/config"
	be "github.com/pickjunk/bgo/error"
	"github.com/pickjunk/zerolog"
	zlog "github.com/pickjunk/zerolog/log"
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

var inner = zlog.With().Str("component", "bgo.log").Logger()
var outer = zlog.With().Logger()

func init() {
	logPath := bc.Get("log").String()
	if logPath != "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			inner.Fatal().Err(err).Send()
		}
		inner.Info().Str("file", logPath).Msg("log redirect")

		outer = outer.Output(f)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		outer = outer.Output(zerolog.ConsoleWriter{Out: os.Stdout})
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
