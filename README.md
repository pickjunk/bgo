# bgo

Business-Go Framework.

### Features

- High-Performance Router base on [httprouter](https://github.com/julienschmidt/httprouter)
- Pretty and neat Middewares
- Uniform logger format. Logger base on [logrus](https://github.com/sirupsen/logrus)
- Support opentracing, integrate jaeger-client
- Support

### Quick Start

- create main.go

```golang
package main

import (
  "net/http"

  dbr "github.com/gocraft/dbr"
  httprouter "github.com/julienschmidt/httprouter"
  bgo "github.com/pickjunk/bgo"
  bgoDbr "github.com/pickjunk/bgo/dbr"
  config "github.com/uber/jaeger-client-go/config"
)

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

  // comment out these code if you don't have jaeger
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

  r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    w.Write([]byte("hello world!"))
  })

  // comment out these code if you don't have mysql
  rWithDbr := r.Middlewares(bgoDbr.Middleware(nil))

  rWithDbr.GET("/dbr", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    ctx := r.Context()
    db := ctx.Value(bgo.CtxKey("dbr")).(*dbr.Session)

    var test struct{}
    err := db.Select(`"empty"`).LoadOneContext(ctx, &test)
    if err != nil {
      bgo.Log.Panic(err)
    }

    w.Write([]byte(`dbr: SELECT "empty"`))
  })

  // comment out these code if you don't like graphql
  r.Graphql("/graphql", g)

  r.ListenAndServe()
}
```

- install dependencies (go modules)

```shell
$ go mod init [your-project-name]
$ go mod tidy
```

- start server

```shell
$ go run main.go
```

### API

### Example

- [[bgo-example]](https://github.com/pickjunk/bgo-example) - A simple project to show how to play with bgo.
