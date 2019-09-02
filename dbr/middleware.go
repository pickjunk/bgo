package dbr

import (
	"context"

	dbr "github.com/gocraft/dbr"
	b "github.com/pickjunk/bgo"
)

// Middleware inject dbr session to context
func Middleware(db *dbr.Session) b.Middleware {
	if db == nil {
		db = New()
	}

	return func(ctx context.Context, next b.Handle) {
		ctx = b.WithValue(ctx, "dbr", db)
		next(ctx)
	}
}

// Dbr get dbr session from context
func Dbr(ctx context.Context) *dbr.Session {
	return b.Value(ctx, "dbr").(*dbr.Session)
}
