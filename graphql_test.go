package bgo

import (
	"testing"
)

func TestMergeSchema(t *testing.T) {
	g := &Graphql{
		schema: "type Query {a(b:B): A} type Mutation {a(b:B): A} type A {ab: ID}",
	}
	g.MergeSchema("type Query {B(a:A): B} type Mutation {a(b:B): A} type B {ab: String}")

	if g.schema != `type Query {
a(b:B):A
B(a:A):B
}
type Mutation {
a(b:B):A
a(b:B):A
}
type A {ab: ID}
type B {ab: String}` {
		t.Error(g.schema)
		t.Error("MergeSchema fail")
	}

	g = &Graphql{
		schema: "type Mutation {a(b:B): A}",
	}
	g.MergeSchema("type Query {B(a:A): B}")

	if g.schema != `type Query {
B(a:A):B
}
type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("MergeSchema fail")
	}

	g = &Graphql{
		schema: "",
	}
	g.MergeSchema("type Query {B(a:A): B} type Mutation {a(b:B): A}")

	if g.schema != `type Query {
B(a:A):B
}
type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("MergeSchema fail")
	}

	g = &Graphql{
		schema: "",
	}
	g.MergeSchema("type Mutation {a(b:B): A}")

	if g.schema != `type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("MergeSchema fail")
	}

	g = &Graphql{
		schema: "type Mutation {a(b:B): A}",
	}
	g.MergeSchema("")

	if g.schema != `type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("MergeSchema fail")
	}
}
