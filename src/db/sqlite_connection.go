package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	Db *sqlx.DB
}

var connection = &Connection{Db: nil}

func GetConnection(dbPath string) (*Connection, error) {
	if connection.Db == nil {
		connection.Db = createConnection(dbPath)
	} else {
		if connection.Db.Ping() != nil {
			connection.Db = createConnection(dbPath)
		}
	}

	return connection, nil
}

func createConnection(dbPath string) *sqlx.DB {
	dsn := "file:" + dbPath + "?cache=shared"
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
