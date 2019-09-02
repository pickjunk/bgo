package bgo

import (
	"context"
	"errors"
	"net/http"

	ot "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	be "github.com/pickjunk/bgo/error"
)

func recoverMiddleware(ctx context.Context, next Handle) {
	defer func() {
		w := Response(ctx)

		if r := recover(); r != nil {
			var err error
			switch t := r.(type) {
			case *be.BusinessError:
				// BusinessError is not an error but a hint
				// just response it here
				w.Write([]byte(t.Error()))
				return
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}

			span := ot.SpanFromContext(ctx)
			if span != nil {
				span.LogFields(otlog.String("event", "error"), otlog.Error(err))
			}

			// log other error, except SystemError
			// SystemError should be log before it raised
			// so do not log it twice here
			if _, ok := err.(*be.SystemError); !ok {
				log.Err(err).Stack().Send()
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	next(ctx)
}
