package dbr

import (
	"context"

	dbr "github.com/gocraft/dbr"
	bgo "github.com/pickjunk/bgo"
)

// Middleware inject dbr session to context
func Middleware(conn *dbr.Connection) bgo.Middleware {
	if conn == nil {
		conn = New()
	}

	return func(ctx context.Context, next bgo.Handle) {
		db := conn.NewSession(nil)
		ctx = context.WithValue(ctx, bgo.CtxKey("dbr"), db)

		next(ctx)
	}
}
