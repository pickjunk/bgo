package bgo

import (
	"context"
	"regexp"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

// Graphql struct
type Graphql struct {
	schema   string
	resolver interface{}
}

type graphqlLogger struct{}

// LogPanic is used to log recovered panic values that occur during query execution
func (l *graphqlLogger) LogPanic(_ context.Context, value interface{}) {
	// log panic error in relay.go, should not log here
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
