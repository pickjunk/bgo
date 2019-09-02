package bgo

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"time"

	ot "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
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
	w := Response(ctx)
	r := Request(ctx)
	ps := Params(ctx)

	span, ctx := ot.StartSpanFromContext(ctx, "http.handle")
	defer span.Finish()

	sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	ctx = withValue(ctx, "http", &HTTP{sw, r, ps})

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
		log.Info().
			Str("ip", ip).
			Str("host", r.Host).
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Int("status", sw.status).
			Int("length", sw.length).
			Str("ua", r.Header.Get("User-Agent")).
			Str("referer", r.Header.Get("Referer")).
			Dur("duration", duration).
			Msg("http handle")
	} else {
		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Int("status", sw.status).
			Dur("duration", duration).
			Msg("http.handle")
	}
}
