package dbr

import (
	"context"

	dbr "github.com/gocraft/dbr"
	bgo "github.com/pickjunk/bgo"
)

// Middleware inject dbr session to context
func Middleware(db *dbr.Session) bgo.Middleware {
	if db == nil {
		db = New()
	}

	return func(ctx context.Context, next bgo.Handle) {
		ctx = context.WithValue(ctx, bgo.CtxKey("dbr"), db)
		next(ctx)
	}
}
