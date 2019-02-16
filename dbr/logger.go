package dbr

import (
	dbr "github.com/gocraft/dbr/opentracing"
	bgo "github.com/pickjunk/bgo"
	log "github.com/sirupsen/logrus"
)

// Logger for dbr
type Logger struct {
	*bgo.Logger
	*dbr.EventReceiver
}

// NewLogger for dbr
func NewLogger() *Logger {
	return &Logger{
		bgo.Log,
		&dbr.EventReceiver{},
	}
}

// ---------- implements dbr EventReceiver interface ----------
// https://github.com/gocraft/dbr/blob/master/event.go

func (l *Logger) kvs2Fields(kvs map[string]string) log.Fields {
	fields := log.Fields{}
	for k, v := range kvs {
		fields[k] = v
	}
	return fields
}

// Event func
func (l *Logger) Event(eventName string) {
	l.Info(eventName)
}

// EventKv func
func (l *Logger) EventKv(eventName string, kvs map[string]string) {
	l.WithFields(l.kvs2Fields(kvs)).Info(eventName)
}

// EventErr func
func (l *Logger) EventErr(eventName string, err error) error {
	l.WithField("msg", err).Error(eventName)
	return err
}

// EventErrKv func
func (l *Logger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	fields := l.kvs2Fields(kvs)
	fields["msg"] = err
	l.WithFields(fields).Error(eventName)
	return err
}

// Timing func
func (l *Logger) Timing(eventName string, nanoseconds int64) {
	l.WithField("timing", nanoseconds).Info(eventName)
}

// TimingKv func
func (l *Logger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	fields := l.kvs2Fields(kvs)
	fields["timing"] = nanoseconds
	l.WithFields(fields).Info(eventName)
}
