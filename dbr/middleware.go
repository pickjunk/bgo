package dbr

import (
	"context"

	b "github.com/pickjunk/bgo"
)

// Middleware inject dbr session to context
func Middleware(db *DB) b.Middleware {
	if db == nil {
		db = New()
	}

	return func(ctx context.Context, next b.Handle) {
		ctx = b.WithValue(ctx, "dbr", db)
		next(ctx)
	}
}

// Dbr get dbr session from context
func Dbr(ctx context.Context) *DB {
	return b.Value(ctx, "dbr").(*DB)
}
