package bgo

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"time"

	ot "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	bl "github.com/pickjunk/bgo/log"
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

func ip(r *http.Request) string {
	// client ip
	ip := r.RemoteAddr
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip = xff
	}
	// clear port
	ip = regexp.MustCompile(`\:\d+$`).ReplaceAllString(ip, "")
	return ip
}

func logMiddleware(ctx context.Context, next Handle) {
	w := Response(ctx)
	r := Request(ctx)
	ps := Params(ctx)

	l := log.With().
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Logger()
	if os.Getenv("ENV") == "production" {
		l = l.With().
			Str("ip", ip(r)).
			Str("host", r.Host).
			Str("ua", r.Header.Get("User-Agent")).
			Str("referer", r.Header.Get("Referer")).
			Logger()
	}
	ctx = withValue(ctx, "log", &bl.Logger{Logger: l})

	sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	ctx = withValue(ctx, "http", &HTTP{sw, r, ps})

	span, ctx := ot.StartSpanFromContext(ctx, "http")
	defer span.Finish()
	start := time.Now()
	next(ctx)
	duration := time.Now().Sub(start)

	otext.HTTPMethod.Set(span, r.Method)
	otext.HTTPUrl.Set(span, r.RequestURI)
	otext.HTTPStatusCode.Set(span, uint16(sw.status))
	if sw.status >= http.StatusInternalServerError {
		otext.Error.Set(span, true)
	}

	l.Info().
		Int("status", sw.status).
		Int("length", sw.length).
		Dur("duration", duration).
		Msg("http")
}

// Log get contextual logger
func Log(ctx context.Context) *bl.Logger {
	return value(ctx, "log").(*bl.Logger)
}
