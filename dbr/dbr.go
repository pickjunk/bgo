package dbr

import (
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr"
	bc "github.com/pickjunk/bgo/config"
)

// DB session instance, base on dbr.Session
type DB struct {
	*dbr.Session
}

// New dbr.Session
// if optionalDSN is omited, mysql.dsn in config.yml will be used
// Supported config(config.yml):
// mysql.dsn - dsn for connection
// mysql.maxIdleConns - max idle connections for pool
// mysql.maxOpenConns - max open connections for pool
func New(optionalDSN ...string) *DB {
	var dsn string

	if len(optionalDSN) > 0 {
		dsn = optionalDSN[0]
	}

	if dsn == "" {
		dsn = bc.Get("mysql.dsn").String()
		if dsn == "" {
			log.Panic().Str("field", "mysql.dsn").Msg("config field not found")
		}
	}

	// open connection
	conn, err := dbr.Open("mysql", dsn, log)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	maxIdleConns := int(bc.Get("mysql.maxIdleConns").Int())
	if maxIdleConns == 0 {
		maxIdleConns = 1
	}
	conn.DB.SetMaxIdleConns(maxIdleConns)

	maxOpenConns := int(bc.Get("mysql.maxOpenConns").Int())
	if maxOpenConns == 0 {
		maxOpenConns = 1
	}
	conn.DB.SetMaxOpenConns(maxOpenConns)

	// ping
	err = conn.DB.Ping()
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().
		Str("dsn", dsn).
		Int("maxIdleConns", maxIdleConns).
		Int("maxOpenConns", maxOpenConns).
		Msg("dbr open")

	return &DB{
		conn.NewSession(nil),
	}
}

// export dbr expression function for convenience
var (
	// And creates AND from a list of conditions.
	And = dbr.And
	// Or creates OR from a list of conditions.
	Or = dbr.Or
	// Eq is `=`.
	// When value is nil, it will be translated to `IS NULL`.
	// When value is a slice, it will be translated to `IN`.
	// Otherwise it will be translated to `=`.
	Eq = dbr.Eq
	// Neq is `!=`.
	// When value is nil, it will be translated to `IS NOT NULL`.
	// When value is a slice, it will be translated to `NOT IN`.
	// Otherwise it will be translated to `!=`.
	Neq = dbr.Neq
	// Gt is `>`.
	Gt = dbr.Gt
	// Gte is '>='.
	Gte = dbr.Gte
	// Lt is '<'.
	Lt = dbr.Lt
	// Lte is `<=`.
	Lte = dbr.Lte
	// Like is `LIKE`, with an optional `ESCAPE` clause
	Like = dbr.Like
	// NotLike is `NOT LIKE`, with an optional `ESCAPE` clause
	NotLike = dbr.NotLike
	// Expr allows raw expression to be used when current SQL syntax is
	// not supported by gocraft/dbr.
	Expr = dbr.Expr
	// Union builds
	Union = dbr.Union
	// UnionAll builds
	UnionAll = dbr.UnionAll
)
