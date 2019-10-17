package log

import (
	"runtime"

	zl "github.com/pickjunk/zerolog"
)

const contextCallerSkipFrameCount = 2

// fork from zerolog
// to support severity-based caller logging
type callerHook struct{}

func (h callerHook) Run(e *zl.Event, level zl.Level, msg string) {
	switch level {
	case zl.ErrorLevel, zl.PanicLevel, zl.FatalLevel:
		_, file, line, ok := runtime.Caller(
			zl.CallerSkipFrameCount + contextCallerSkipFrameCount,
		)
		if ok {
			e.Str(zl.CallerFieldName, zl.CallerMarshalFunc(file, line))
		}
	}
}
