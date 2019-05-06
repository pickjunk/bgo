package bgo

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"time"

	ot "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
)

// https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func logMiddleware(ctx context.Context, next Handle) {
	httpCtx := ctx.Value(CtxKey("http")).(*HTTP)
	w := httpCtx.Response
	r := httpCtx.Request
	ps := httpCtx.Params

	span, ctx := ot.StartSpanFromContext(ctx, "http.handle")
	defer span.Finish()

	sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	ctx = context.WithValue(ctx, CtxKey("http"), &HTTP{sw, r, ps})

	start := time.Now()
	next(ctx)
	duration := time.Now().Sub(start)

	otext.HTTPMethod.Set(span, r.Method)
	otext.HTTPUrl.Set(span, r.RequestURI)
	otext.HTTPStatusCode.Set(span, uint16(sw.status))

	if sw.status >= http.StatusInternalServerError {
		otext.Error.Set(span, true)
	}

	// client ip
	ip := r.RemoteAddr
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip = xff
	}
	// clear port
	re := regexp.MustCompile(`\:\d+$`)
	ip = re.ReplaceAllString(ip, "")

	if os.Getenv("ENV") == "production" {
		Log.WithFields(log.Fields{
			"ip":       ip,
			"host":     r.Host,
			"method":   r.Method,
			"uri":      r.RequestURI,
			"status":   sw.status,
			"length":   sw.length,
			"ua":       r.Header.Get("User-Agent"),
			"duration": duration,
		}).Info("http.handle")
	} else {
		Log.WithFields(log.Fields{
			"method":   r.Method,
			"uri":      r.RequestURI,
			"status":   sw.status,
			"duration": duration,
		}).Info("http.handle")
	}
}
