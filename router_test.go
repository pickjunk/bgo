package bgo

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	cors "github.com/rs/cors"
	assert "github.com/stretchr/testify/assert"
)

func TestMiddlewares(t *testing.T) {
	r := New()
	foo := 0

	r.Middlewares(
		func(ctx context.Context, next Handle) {
			foo++
			next(ctx)

			if foo != 0 {
				t.Error("fail")
			}

			foo = 200
		},
		func(ctx context.Context, next Handle) {
			if foo != 1 {
				t.Error("fail")
			}
			foo++

			next(ctx)

			if foo != 1 {
				t.Error("fail")
			}
			foo--
		},
	).GET("/", func(ctx context.Context) {
		if foo != 2 {
			t.Error("fail")
		}
		foo--
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, 200, foo)

	foo = 0
	r.Middlewares(
		func(ctx context.Context, next Handle) {
			foo++
			next(ctx)

			if foo != 0 {
				t.Error("fail")
			}

			foo = 200
		},
	).POST("/", func(ctx context.Context) {
		if foo != 1 {
			t.Error("fail")
		}
		foo--
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil))
	assert.Equal(t, 200, foo)
}

func TestHandle(t *testing.T) {
	r := New()
	foo := 0

	r.GET(
		"/",
		func(ctx context.Context, next Handle) {
			foo++
			next(ctx)

			if foo != 0 {
				t.Error("fail")
			}
			foo = 200
		},
		func(ctx context.Context, next Handle) {
			if foo != 1 {
				t.Error("fail")
			}
			foo++

			next(ctx)

			if foo != 1 {
				t.Error("fail")
			}
			foo--
		},
		func(ctx context.Context) {
			if foo != 2 {
				t.Error("fail")
			}
			foo--
		},
	)
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, 200, foo)

	foo = 0
	r.POST(
		"/",
		func(ctx context.Context, next Handle) {
			foo++
			next(ctx)

			if foo != 0 {
				t.Error("fail")
			}
			foo = 200
		},
		func(ctx context.Context) {
			if foo != 1 {
				t.Error("fail")
			}
			foo--
		},
	)
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil))
	assert.Equal(t, 200, foo)

	// block log print, avoid panic msg to print
	Log.SetOutput(nil)
	defer Log.SetOutput(os.Stdout)

	assert.Panics(t, func() {
		r.GET(
			"/",
		)
	})
	assert.Panics(t, func() {
		r.GET(
			"/",
			func(ctx context.Context, next Handle) {
			},
		)
	})
	assert.Panics(t, func() {
		r.GET(
			"/",
			func(ctx context.Context) {
			},
			func(ctx context.Context) {
			},
		)
	})
}

func TestPrefix(t *testing.T) {
	r := New()
	foo := 0

	r.Prefix("/prefix").
		GET("/", func(ctx context.Context) {
			foo++
		}).
		POST("/", func(ctx context.Context) {
			foo++
		}).
		Prefix("/prefix").
		GET("/", func(ctx context.Context) {
			foo++
		})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/prefix/", nil))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/prefix/", nil))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/prefix/prefix/", nil))
	assert.Equal(t, 3, foo)

	r.Prefix("/ab").GET("c", func(ctx context.Context) {
		foo++
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/abc", nil))
	assert.Equal(t, 4, foo)
}

func TestCORS(t *testing.T) {
	r := New()
	foo := 0

	r.CORS(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		// Debug:            true,
	})).POST("/cors", func(ctx context.Context) {
		foo++
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/cors", nil)
	req.Header.Add("Origin", "foo")
	r.ServeHTTP(w, req)
	assert.Equal(t, "Origin", w.Header().Get("Vary"))

	w = httptest.NewRecorder()
	req = httptest.NewRequest("OPTIONS", "/cors", nil)
	req.Header.Add("Origin", "foo")
	req.Header.Add("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)
	assert.Equal(t, "Origin", w.Header().Get("Vary"))

	r.GET("/_cors", func(ctx context.Context) {
		foo++
	})

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/_cors", nil))
	assert.Equal(t, "", w.Header().Get("Vary"))

	assert.Equal(t, 2, foo)
}
