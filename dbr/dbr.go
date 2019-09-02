package dbr

import (
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr"
	b "github.com/pickjunk/bgo"
)

// New dbr.Session
// if optionalDSN is omited, mysql.dsn in config.yml will be used
// Supported config(config.yml):
// mysql.dsn - dsn for connection
// mysql.maxIdleConns - max idle connections for pool
// mysql.maxOpenConns - max open connections for pool
func New(optionalDSN ...string) *dbr.Session {
	var dsn string

	if len(optionalDSN) > 0 {
		dsn = optionalDSN[0]
	}

	if dsn == "" {
		dsn = b.Config.Get("mysql.dsn").String()
		if dsn == "" {
			log.Panic().Str("field", "mysql.dsn").Msg("config field not found")
		}
	}

	// open connection
	conn, err := dbr.Open("mysql", dsn, log)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	maxIdleConns := int(b.Config.Get("mysql.maxIdleConns").Int())
	if maxIdleConns == 0 {
		maxIdleConns = 1
	}
	conn.DB.SetMaxIdleConns(maxIdleConns)

	maxOpenConns := int(b.Config.Get("mysql.maxOpenConns").Int())
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

	return conn.NewSession(nil)
}
