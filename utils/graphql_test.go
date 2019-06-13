package utils

import (
	"context"
	"net/http"
	"regexp"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	req "github.com/imroc/req"
	b "github.com/pickjunk/bgo"
)

type resolver struct{}

var g = b.NewGraphql(&resolver{})

func init() {
	g.MergeSchema(`
	type Query {
		test(one: Int, two: String): Test
	}

	type Test {
		id: String
	}`)
}

type Test struct {
	ID string
}

type TestResolver struct {
	t *Test
}

// Test resolver
func (r *resolver) Test(
	ctx context.Context,
	args struct {
		One *int32
		Two *string
	},
) *TestResolver {
	if args.One != nil && *args.One == 1 {
		b.Throw(100, "test error")
	}

	h := ctx.Value(b.CtxKey("http")).(*b.HTTP)
	w := h.Response
	hr := h.Request

	auth := hr.Header.Get("Authorization")
	m := regexp.MustCompile(`^Bearer (.*)`).FindStringSubmatch(auth)
	if m == nil || m[1] == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	return &TestResolver{
		t: &Test{
			ID: m[1],
		},
	}
}

func (r *TestResolver) ID() *string {
	return &r.t.ID
}

func setup() {
	go func() {
		r := b.New()
		r.Graphql("/api", g)
		r.ListenAndServe()
	}()

	time.Sleep(time.Duration(3) * time.Second)
}

func TestGraphql(t *testing.T) {
	setup()

	var result struct {
		D struct {
			T struct {
				ID string `json:"id"`
			} `json:"test"`
		} `json:"data"`
	}

	r := Graphql{
		URL: "http://localhost:8888/api",
		Query: `
		query ($one: Int, $two: String) {
			test(one: $one, two: $two) {
				id
			}
		}
		`,
	}
	err := r.Fetch(&result)
	if err == nil || err.Error() != "http status error: 401" {
		t.Errorf("can not catch http status error correctly")
	}

	r = Graphql{
		URL: "http://localhost:8888/api",
		Query: `
		query ($one: Int, $two: String) {
			test(one: $one, two: $two) {
				id
			}
		}
		`,
		Variables: map[string]interface{}{
			"one": 1,
		},
		Headers: req.Header{
			"Authorization": "Bearer 123",
		},
	}
	err = r.Fetch(&result)
	_, ok := err.(*b.BusinessError)
	if !ok {
		t.Errorf("can not catch business error correctly")
	}

	r = Graphql{
		URL: "http://localhost:8888/api",
		Query: `
		query ($one: Int, $two: String) {
			test(one: $one, two: $two) {
				id
			}
		}
		`,
		Headers: req.Header{
			"Authorization": "Bearer 123",
		},
	}
	err = r.Fetch(&result)
	if err != nil {
		t.Errorf("can not fetch correctly")
	}

	if result.D.T.ID != "123" {
		t.Errorf("result not correct")
	}
}
