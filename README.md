# bgo

Business-Go Framework.

### Features

- High-Performance Router base on [httprouter](https://github.com/julienschmidt/httprouter)
- Pretty Middewares
- Uniform logger format. Logger base on [logrus](https://github.com/sirupsen/logrus)
- Opentracing, integrate jaeger-client
- Graphql, base on [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)

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

### API

### Example

- [[bgo-example]](https://github.com/pickjunk/bgo-example) - An example to show how to play with bgo.
