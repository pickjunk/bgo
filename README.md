# bgo
[![GoDoc](https://godoc.org/github.com/pickjunk/bgo?status.svg)](https://godoc.org/github.com/pickjunk/bgo)

Business-Go Framework.

### Features

- High-performance router, base on [httprouter](https://github.com/julienschmidt/httprouter)
- Pretty middewares
- Uniform logger, base on [logrus](https://github.com/sirupsen/logrus)
- Opentracing, integrate jaeger-client
- Graphql, base on [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)
- CORS, base on [rs/cors](https://github.com/rs/cors)

### Quick Start

```golang
// main.go
package main

import (
  "net/http"

  dbr "github.com/gocraft/dbr"
  httprouter "github.com/julienschmidt/httprouter"
  bgo "github.com/pickjunk/bgo"
  bgoDbr "github.com/pickjunk/bgo/dbr"
)

func main() {
  r := bgo.New()

  r.GET("/:name", func(ctx context.Context) {
    h := ctx.Value(bgo.CtxKey("http")).(*bgo.HTTP)
    w := h.Response // http.ResponseWriter
    ps := h.Params // httprouter.Params

    w.Write([]byte("hello "+ps.ByName("name")+"!"))
  })

  // comment out these code if you don't use dbr
  rWithDbr := r.Prefix("/dbr").Middlewares(bgoDbr.Middleware(nil))

  rWithDbr.GET("/empty", func(ctx context.Context) {
    h := ctx.Value(bgo.CtxKey("http")).(*bgo.HTTP)
    w := h.Response // http.ResponseWriter
    db := ctx.Value(bgo.CtxKey("dbr")).(*dbr.Session)

    var test struct{}
    err := db.Select(`"empty"`).LoadOneContext(ctx, &test)
    if err != nil {
      bgo.Log.Panic(err)
    }

    w.Write([]byte(`dbr: SELECT "empty"`))
  })

  r.ListenAndServe()
}
```

### HTTP Context

```golang
r.GET("/:name", func(ctx context.Context) {
  h := ctx.Value(bgo.CtxKey("http")).(*bgo.HTTP)
  r := h.Request // *http.Request
  w := h.Response // http.ResponseWriter
  ps := h.Params // httprouter.Params
})
```

### Middlewares

```golang
r.Middlewares(
  func(ctx context.Context, next bgo.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context, next bgo.Handle) {
    // do something

    next(ctx)
  },
).GET(
  "/"
  func(ctx context.Context, next bgo.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context) {
    // do something
  },
)
```

### SubRoute (Prefix + Middlewares)

```golang
subRoute1 := r.Prefix("/sub1")

subRoute2 := r.Prefix("/sub1").Middlewares(
  func(ctx context.Context, next bgo.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context, next bgo.Handle) {
    // do something

    next(ctx)
  },
)
```

### Graphql

```golang
type resolver struct{}

var g = bgo.NewGraphql(&resolver{})

func init() {
  g.MergeSchema(`
  schema {
    query: Query
  }

  type Query {
    greeting(name: String): String!
  }
  `)
}

func (r *resolver) Greeting(
	ctx context.Context,
	args struct {
		Name   *string
	},
) string {
  if args.Name == nil {
    name := "world"
    args.Name = &name
  }
	return "hello " + args.Name
}

func main() {
  r := bgo.New()

  r.Graphql("/graphql", g)
}
```

### Opentracing (jaeger-client)

```golang
package main

import (
  bgo "github.com/pickjunk/bgo"
  config "github.com/uber/jaeger-client-go/config"
)

func main() {
  r := bgo.New()

  closer := bgo.Jaeger(&config.Configuration{
    ServiceName: "bgo-example",
    Sampler: &config.SamplerConfig{
      Type:  "const",
      Param: 1,
    },
    Reporter: &config.ReporterConfig{
      LogSpans: true,
    },
  })
  defer closer.Close()

  r.ListenAndServe()
}
```

### CORS

```golang
package main

import (
  bgo "github.com/pickjunk/bgo"
  "github.com/rs/cors"
)

func main() {
  r := bgo.New()

  r.CORS(cors.AllowAll()).GET("/", func(ctx context.Context) {
    // do something
  })

  r.ListenAndServe()
}
```

### Logger

```golang
package main

import (
  bgo "github.com/pickjunk/bgo"
)

func main() {
  logger := bgo.Log // logger is a *logrus.Logger

  // log everything as logrus
  // https://github.com/sirupsen/logrus
}
```

### Business Error

```golang
package main

import (
  bgo "github.com/pickjunk/bgo"
)

func main() {
  r := bgo.New()

  r.GET("/", func(ctx context.Context) {
    // Throw will trigger a panic, which internal recover middleware
    // will catch and unmarshal to the response content as
    // `{"code":10001,"msg":"passwd error"}`
    bgo.Throw(10001, "passwd error")
  })

  r.ListenAndServe()
}
```

### Example

- [[bgo-example]](https://github.com/pickjunk/bgo-example) - An example to show how to play with bgo.
