package persistence

import (
	"database/sql"
)

var (
	db       *sql.DB
	host     string
	port     string
	user     string
	password string
	database string
	sslmode  string
)

func InitDB(pgHost, pgPort, pgUser, pgPassword, pgDatabase, pgSSLMode string) (*sql.DB, error) {
	host = pgHost
	port = pgPort
	user = pgUser
	password = pgPassword
	database = pgDatabase
	sslmode = pgSSLMode

	return GetDBSession()
}

func SetDB(database *sql.DB) {
	db = database
}
