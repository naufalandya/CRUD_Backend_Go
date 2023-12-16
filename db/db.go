package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenDB() (*sql.DB, error) {
	conn := "postgresql://postgres:Phoeblex25@127.0.0.1/MyZoo?sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
