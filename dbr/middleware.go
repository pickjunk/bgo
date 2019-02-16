package dbr

import (
	"context"
	"net/http"

	bgo "github.com/ChieveiT/bgo"
	dbr "github.com/gocraft/dbr"
	httprouter "github.com/julienschmidt/httprouter"
)

// DbrMiddleware inject dbr session to context
func DbrMiddleware(conn *dbr.Connection) bgo.Middleware {
	if conn == nil {
		conn = New()
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, next httprouter.Handle) {
		ctx := r.Context()
		db := conn.NewSession(nil)
		ctx = context.WithValue(ctx, bgo.CtxKey("dbr"), db)

		next(w, r.WithContext(ctx), ps)
	}
}
