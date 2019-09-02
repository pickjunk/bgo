package bgo

import (
	"context"
	"os"
	"regexp"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	be "github.com/pickjunk/bgo/error"
)

// Graphql struct
type Graphql struct {
	schema   string
	resolver interface{}
}

type graphqlLogger struct{}

// LogPanic is used to log recovered panic values that occur during query execution
func (l *graphqlLogger) LogPanic(_ context.Context, value interface{}) {
	// skip business error
	if _, ok := value.(*be.BusinessError); ok {
		return
	}

	// skip system error
	if _, ok := value.(*be.SystemError); ok {
		return
	}

	log.Error().Msgf("graphql: %v", value)
}

// Graphql create a graphql endpoint
func (r *Router) Graphql(path string, g *Graphql) *Router {
	schema := graphql.MustParseSchema(
		g.schema,
		g.resolver,
		graphql.Logger(&graphqlLogger{}),
	)

	r.GET(path, func(ctx context.Context) {
		relay(ctx, schema)
	})
	r.POST(path, func(ctx context.Context) {
		relay(ctx, schema)
	})

	if os.Getenv("ENV") != "production" {
		r.GET(path+"-i", func(ctx context.Context) {
			h := ctx.Value(CtxKey("http")).(*HTTP)
			h.Response.Write([]byte(g.Graphiql(path)))
		})
	}

	return r
}

// NewGraphql create a Graphql struct
func NewGraphql(resolver interface{}) *Graphql {
	return &Graphql{
		schema: `
		schema {
			query: Query
			mutation: Mutation
		}

		type Query {}
		type Mutation {}
		`,
		resolver: resolver,
	}
}

func (g *Graphql) replaceFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

// MergeSchema merge two graphql schema
func (g *Graphql) MergeSchema(s string) {
	a := g.schema
	b := s

	r := regexp.MustCompile("(?s)type Query {(.*?)}")
	var query []string
	a = g.replaceFunc(r, a, func(m []string) string {
		lines := strings.Split(strings.Replace(m[1], " ", "", -1), "\n")
		query = append(query, lines...)
		return ""
	})
	b = g.replaceFunc(r, b, func(m []string) string {
		lines := strings.Split(strings.Replace(m[1], " ", "", -1), "\n")
		query = append(query, lines...)
		return ""
	})

	r = regexp.MustCompile("(?s)type Mutation {(.*?)}")
	var mutation []string
	a = g.replaceFunc(r, a, func(m []string) string {
		lines := strings.Split(strings.Replace(m[1], " ", "", -1), "\n")
		mutation = append(mutation, lines...)
		return ""
	})
	b = g.replaceFunc(r, b, func(m []string) string {
		lines := strings.Split(strings.Replace(m[1], " ", "", -1), "\n")
		mutation = append(mutation, lines...)
		return ""
	})

	result := ""
	if len(query) > 0 {
		result += "type Query {\n" + strings.Join(query, "\n") + "\n}\n"
	}
	if len(mutation) > 0 {
		result += "type Mutation {\n" + strings.Join(mutation, "\n") + "\n}\n"
	}
	result += strings.Trim(a, "\n ") + "\n" + strings.Trim(b, "\n ")
	result = strings.Trim(result, "\n ")

	g.schema = result
}

// Graphiql html
func (g *Graphql) Graphiql(uri string) string {
	return `
		<!DOCTYPE html>
		<html>
			<head>
				<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
				<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
				<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
				<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
				<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
			</head>
			<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
				<div id="graphiql" style="height: 100vh;">Loading...</div>
				<script>
					function graphQLFetcher(graphQLParams) {
						return fetch("` + uri + `", {
							method: "post",
							body: JSON.stringify(graphQLParams),
							credentials: "include",
						}).then(function (response) {
							return response.text();
						}).then(function (responseBody) {
							try {
								return JSON.parse(responseBody);
							} catch (error) {
								return responseBody;
							}
						});
					}
					ReactDOM.render(
						React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
						document.getElementById("graphiql")
					);
				</script>
			</body>
		</html>
	`
}
