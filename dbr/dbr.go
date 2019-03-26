package dbr

import (
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr"
	bgo "github.com/pickjunk/bgo"
)

// New open dbr.Connection with config.yml
func New() *dbr.Connection {
	config, ok := bgo.Config["mysql"].(map[interface{}]interface{})
	if !ok {
		bgo.Log.Panic("mysql config not found")
	}

	dsn, ok := config["dsn"].(string)
	if !ok {
		bgo.Log.Panic("mysql dsn is required")
	}

	conn, err := dbr.Open("mysql", dsn, NewLogger())
	if err != nil {
		bgo.Log.Panic(err)
	}

	maxIdleConns, ok := config["maxIdleConns"].(int)
	if !ok {
		maxIdleConns = 5
	}
	conn.DB.SetMaxIdleConns(maxIdleConns)
	maxOpenConns, ok := config["maxOpenConns"].(int)
	if !ok {
		maxOpenConns = 10
	}
	conn.DB.SetMaxOpenConns(maxOpenConns)

	err = conn.DB.Ping()
	if err != nil {
		bgo.Log.Panic(err)
	}

	bgo.Log.WithField("dsn", dsn).Info("dbr.Open")

	return conn
}

// Simple create db instance without config.yml for quick usage
// only one conn will be created, conn pool will be disabled
func Simple(dsn string) *dbr.Session {
	conn, err := dbr.Open("mysql", dsn, NewLogger())
	if err != nil {
		bgo.Log.Panic(err)
	}

	conn.DB.SetMaxIdleConns(1)
	conn.DB.SetMaxOpenConns(1)

	err = conn.DB.Ping()
	if err != nil {
		bgo.Log.Panic(err)
	}

	bgo.Log.WithField("dsn", dsn).Info("dbr.Open")

	return conn.NewSession(nil)
}
