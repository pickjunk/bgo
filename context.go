package bgo

import (
	"context"
)

type innerKey string

func (c innerKey) String() string {
	return "bgo inner context key: " + string(c)
}

func withValue(ctx context.Context, key string, v interface{}) context.Context {
	return context.WithValue(ctx, innerKey(key), v)
}

func value(ctx context.Context, key string) interface{} {
	return ctx.Value(innerKey(key))
}

type outerKey string

func (c outerKey) String() string {
	return "bgo outer context key: " + string(c)
}

// WithValue create a new context with a specific key & value
func WithValue(ctx context.Context, key string, v interface{}) context.Context {
	return context.WithValue(ctx, outerKey(key), v)
}

// Value get value with a specific key
func Value(ctx context.Context, key string) interface{} {
	return ctx.Value(outerKey(key))
}
