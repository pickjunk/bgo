package bgo

import (
	"errors"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
	ot "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

func recoverMiddleware(w http.ResponseWriter, req *http.Request, ps httprouter.Params, next httprouter.Handle) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch t := r.(type) {
			case *BusinessError:
				w.Write([]byte(t.Error()))
				return
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}

			ctx := req.Context()
			span := ot.SpanFromContext(ctx)
			if span != nil {
				span.LogFields(otlog.String("event", "error"), otlog.Error(err))
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	next(w, req, ps)
}
