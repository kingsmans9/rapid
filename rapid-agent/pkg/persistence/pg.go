//go:build !testing

package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetDBSession() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, database, sslmode)
	newDB, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	db = newDB
	return db, nil
}
