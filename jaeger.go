package bgo

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	config "github.com/uber/jaeger-client-go/config"
)

type jaegerLogger struct{}

func (l *jaegerLogger) Error(msg string) {
	log.Error().Str("component", "bgo.jaeger").Msg(msg)
}

// Infof logs a message at info priority
func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	log.Debug().Str("component", "bgo.jaeger").Msgf(msg, args...)
}

// Jaeger setup a jaeger tracer
func Jaeger(cfg *config.Configuration) io.Closer {
	tracer, closer, err := cfg.NewTracer(config.Logger(&jaegerLogger{}))
	if err != nil {
		log.Panic().Err(err).Send()
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}
