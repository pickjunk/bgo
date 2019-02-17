package dbr

import (
	dbr "github.com/gocraft/dbr"
	bgo "github.com/pickjunk/bgo"
)

// New dbr.Connection
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

	bgo.Log.Info("dbr connected")

	return conn
}
