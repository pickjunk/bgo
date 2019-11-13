package dbr

import (
	"time"
	"errors"

	dbr "github.com/gocraft/dbr/opentracing"
	bl "github.com/pickjunk/bgo/log"
)

// Logger for dbr
type Logger struct {
	*bl.Logger
	*dbr.EventReceiver
}

var log = &Logger{
	bl.New("dbr"),
	&dbr.EventReceiver{},
}

// ---------- implements dbr EventReceiver interface ----------
// https://github.com/gocraft/dbr/blob/master/event.go

// Event func
func (l *Logger) Event(eventName string) {
	l.Info().Msg(eventName)
}

// EventKv func
func (l *Logger) EventKv(eventName string, kvs map[string]string) {
	info := l.Info()
	for k, v := range kvs {
		info = info.Str(k, v)
	}
	info.Msg(eventName)
}

// EventErr func
func (l *Logger) EventErr(eventName string, err error) error {
	return errors.New(eventName + ": msg=" + err.Error())
}

// EventErrKv func
func (l *Logger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	msg := eventName + ": msg=" + err.Error()
	for k, v := range kvs {
		msg += " " + k + "=" + v
	}
	return errors.New(msg)
}

// Timing func
func (l *Logger) Timing(eventName string, nanoseconds int64) {
	l.Info().Dur("duration", time.Duration(nanoseconds)*time.Nanosecond).Msg(eventName)
}

// TimingKv func
func (l *Logger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	info := l.Info()
	for k, v := range kvs {
		info = info.Str(k, v)
	}
	info.Dur("duration", time.Duration(nanoseconds)*time.Nanosecond).Msg(eventName)
}
