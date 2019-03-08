package bgo

import (
	"context"
	"errors"
	"net/http"

	ot "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	log "github.com/sirupsen/logrus"
)

func recoverMiddleware(ctx context.Context, next Handle) {
	defer func() {
		httpCtx := ctx.Value(CtxKey("http")).(*HTTP)
		w := httpCtx.Response

		if r := recover(); r != nil {
			var err error
			switch t := r.(type) {
			case *BusinessError:
				w.Write([]byte(t.Error()))
				return
			case *log.Entry:
				// error from Log.Panic, skip logging
				// because Log.Panic has logged the error when it called
				err = nil
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

			if err != nil {
				Log.Error(err)
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	next(ctx)
}
