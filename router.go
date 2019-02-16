package bgo

import (
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
	sentry "github.com/onrik/logrus/sentry"
	cors "github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Router thin wrapper for httprouter.Router
type Router struct {
	prefix      string
	middlewares []Middleware
	cors        *cors.Cors
	*httprouter.Router
}

// New create a bgo Router
func New() *Router {
	// sentry
	if sentryDSN, ok := Config["sentry"].(string); ok {
		sentryHook := sentry.NewHook(sentryDSN, log.ErrorLevel, log.PanicLevel, log.FatalLevel)
		Log.AddHook(sentryHook)
	}

	return &Router{
		Router: httprouter.New(),
		middlewares: []Middleware{
			logMiddleware,
			recoverMiddleware,
		},
	}
}

// Prefix append prefix
func (r *Router) Prefix(p string) *Router {
	new := *r
	new.prefix = new.prefix + p
	return &new
}

// Middlewares register middlewares
func (r *Router) Middlewares(layers ...Middleware) *Router {
	new := *r
	new.middlewares = append(r.middlewares, layers...)
	return &new
}

// CORS register middlewares
func (r *Router) CORS(c *cors.Cors) *Router {
	new := *r
	new.cors = c
	return &new
}

// Handle func
type Handle = func(http.ResponseWriter, *http.Request, httprouter.Params)

// Middleware func
type Middleware = func(http.ResponseWriter, *http.Request, httprouter.Params, httprouter.Handle)

// Handle define a route
func (r *Router) Handle(method, path string, middlewaresAndHandle ...interface{}) *Router {
	l := len(middlewaresAndHandle)
	if l == 0 {
		Log.Panic("expect bgo.Handle")
	}

	finalHandle, ok := middlewaresAndHandle[l-1].(Handle)
	if !ok {
		Log.Panicf("expect bgo.Handle, but get %T", middlewaresAndHandle[l-1])
	}

	middlewares := r.middlewares
	for i := 0; i < l-1; i++ {
		middleware, ok := middlewaresAndHandle[i].(Middleware)
		if !ok {
			Log.Panicf("expect bgo.Middleware, but get %T", middlewaresAndHandle[i])
		}
		middlewares = append(middlewares, middleware)
	}

	// handle wrapped with middlewares
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		next := finalHandle
		finalHandle = func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			middleware(w, req, ps, next)
		}
	}

	hr := r.Router

	// compatible with rs/cors
	if r.cors != nil {
		innerHandle := finalHandle
		finalHandle = func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			r.cors.HandlerFunc(w, req)
			innerHandle(w, req, ps)
		}

		if method != "OPTIONS" {
			hr.Handle("OPTIONS", r.prefix+path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
				r.cors.HandlerFunc(w, req)
			})
		}
	}

	hr.Handle(method, r.prefix+path, finalHandle)

	return r
}

// GET is a shortcut for router.Handle("GET", path, middlewaresAndHandle...)
func (r *Router) GET(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("GET", path, middlewaresAndHandle...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, middlewaresAndHandle...)
func (r *Router) HEAD(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("HEAD", path, middlewaresAndHandle...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, middlewaresAndHandle...)
func (r *Router) OPTIONS(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("OPTIONS", path, middlewaresAndHandle...)
}

// POST is a shortcut for router.Handle("POST", path, middlewaresAndHandle...)
func (r *Router) POST(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("POST", path, middlewaresAndHandle...)
}

// PUT is a shortcut for router.Handle("PUT", path, middlewaresAndHandle...)
func (r *Router) PUT(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("PUT", path, middlewaresAndHandle...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, middlewaresAndHandle...)
func (r *Router) PATCH(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("PATCH", path, middlewaresAndHandle...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, middlewaresAndHandle...)
func (r *Router) DELETE(path string, middlewaresAndHandle ...interface{}) *Router {
	return r.Handle("DELETE", path, middlewaresAndHandle...)
}
