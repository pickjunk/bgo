package dbr

import (
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr"
	bgo "github.com/pickjunk/bgo"
)

// New dbr.Session
// if optionalDSN is omited, mysql.dsn in config.yml will be used
// Supported config(config.yml):
// mysql.dsn - dsn for connection
// mysql.maxIdleConns - max idle connections for pool
// mysql.maxOpenConns - max open connections for pool
func New(optionalDSN ...string) *dbr.Session {
	var dsn string
	var maxIdleConns, maxOpenConns int

	if len(optionalDSN) > 0 {
		dsn = optionalDSN[0]
	}

	config, ok := bgo.Config["mysql"].(map[interface{}]interface{})
	if !ok {
		config = make(map[interface{}]interface{})
	}

	if dsn == "" {
		dsn, ok = config["dsn"].(string)
		if !ok {
			bgo.Log.Panic("mysql dsn is required")
		}
	}

	// open connection
	conn, err := dbr.Open("mysql", dsn, NewLogger())
	if err != nil {
		bgo.Log.Panic(err)
	}

	// connection pool config
	maxIdleConns, ok = config["maxIdleConns"].(int)
	if !ok {
		maxIdleConns = 1
	}
	conn.DB.SetMaxIdleConns(maxIdleConns)
	maxOpenConns, ok = config["maxOpenConns"].(int)
	if !ok {
		maxOpenConns = 1
	}
	conn.DB.SetMaxOpenConns(maxOpenConns)

	// ping
	err = conn.DB.Ping()
	if err != nil {
		bgo.Log.Panic(err)
	}

	bgo.Log.
		WithField("dsn", dsn).
		WithField("maxIdleConns", maxIdleConns).
		WithField("maxOpenConns", maxOpenConns).
		Info("dbr.Open")

	return conn.NewSession(nil)
}
