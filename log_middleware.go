package bgo

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"time"

	ot "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog"
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

	access := make(map[string]string)
	access["method"] = r.Method
	access["uri"] = r.RequestURI
	if os.Getenv("ENV") == "production" {
		access["ip"] = ip(r)
		access["host"] = r.Host
		access["ua"] = r.Header.Get("User-Agent")
		access["referer"] = r.Header.Get("Referer")
	}
	ctx = withValue(ctx, "access", access)

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

	var e *zerolog.Event
	switch s := sw.status; {
	case s >= 200 && s < 300:
		e = log.Info()
	case s >= 300 && s < 500:
		e = log.Warn()
	default:
		e = log.Error()
	}
	for k, v := range access {
		e.Str(k, v)
	}
	e.Int("status", sw.status).
		Int("length", sw.length).
		Dur("duration", duration).
		Msg("access")
}

// Access context, everything added to this map
// will be logged as the field of access log
func Access(ctx context.Context) map[string]string {
	return value(ctx, "access").(map[string]string)
}
