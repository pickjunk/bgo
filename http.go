package bgo

import (
	"context"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
)

// HTTP context
type HTTP struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   httprouter.Params
}

// Request get request from context
func Request(ctx context.Context) *http.Request {
	return value(ctx, "http").(*HTTP).Request
}

// Response get response from context
func Response(ctx context.Context) http.ResponseWriter {
	return value(ctx, "http").(*HTTP).Response
}

// Params get params from context
func Params(ctx context.Context) httprouter.Params {
	return value(ctx, "http").(*HTTP).Params
}

// Param get param from context
func Param(ctx context.Context, key string) string {
	return value(ctx, "http").(*HTTP).Params.ByName(key)
}
