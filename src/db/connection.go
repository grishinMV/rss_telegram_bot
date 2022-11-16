package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

type Connection struct {
	Db *sqlx.DB
}

var connection = &Connection{Db: nil}

func GetConnection(config *Config) (*Connection, error) {
	if connection.Db == nil {
		connection.Db = createConnection(config)
	} else {
		if connection.Db.Ping() != nil {
			connection.Db = createConnection(config)
		}
	}

	return connection, nil
}

func createConnection(config *Config) *sqlx.DB {
	dsn := config.User() + ":" + config.Password() + "@tcp(" + config.Host() + ":" + config.Port() + ")/" + config.DbName() + "?charset=" + config.Charset() + "&parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func LoadConfigFromEnv() *Config {
	host, exist := os.LookupEnv("RSS_DB_HOST")
	if !exist {
		host = "localhost"
	}

	port, exist := os.LookupEnv("RSS_DB_PORT")
	if !exist {
		port = "3306"
	}

	user, exist := os.LookupEnv("RSS_DB_USER")
	if !exist {
		user = "root"
	}

	pass, exist := os.LookupEnv("RSS_DB_PASSWORD")
	if !exist {
		pass = "test"
	}

	charset, exist := os.LookupEnv("RSS_DB_CHARSET")
	if !exist {
		charset = "utf8mb4"
	}

	dbName, exist := os.LookupEnv("RSS_DB_NAME")
	if !exist {
		dbName = "rss_parser"
	}

	return &Config{
		host:     host,
		port:     port,
		user:     user,
		password: pass,
		charset:  charset,
		dbName:   dbName,
	}
}
