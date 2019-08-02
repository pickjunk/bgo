package dbr

import (
	"time"

	dbr "github.com/gocraft/dbr/opentracing"
	b "github.com/pickjunk/bgo"
	logrus "github.com/sirupsen/logrus"
)

// Logger for dbr
type Logger struct {
	*logrus.Entry
	*dbr.EventReceiver
}

// Log log dbr
var log = &Logger{
	b.Log.WithField("prefix", "dbr"),
	&dbr.EventReceiver{},
}

// ---------- implements dbr EventReceiver interface ----------
// https://github.com/gocraft/dbr/blob/master/event.go

func (l *Logger) kvs2Fields(kvs map[string]string) logrus.Fields {
	fields := logrus.Fields{}
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
	l.WithField("duration", (time.Duration(nanoseconds) * time.Nanosecond).String()).Info(eventName)
}

// TimingKv func
func (l *Logger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	fields := l.kvs2Fields(kvs)
	fields["duration"] = (time.Duration(nanoseconds) * time.Nanosecond).String()
	l.WithFields(fields).Info(eventName)
}
