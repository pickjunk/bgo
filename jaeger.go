package bgo

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	config "github.com/uber/jaeger-client-go/config"
)

type jaegerLogger struct{}

func (l *jaegerLogger) Error(msg string) {
	Log.Error("jaeger: " + msg)
}

// Infof logs a message at info priority
func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	Log.Infof("jaeger: "+msg, args...)
}

// Jaeger setup a jaeger tracer
func Jaeger(cfg *config.Configuration) io.Closer {
	tracer, closer, err := cfg.NewTracer(config.Logger(&jaegerLogger{}))
	if err != nil {
		Log.Panic(err)
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}
