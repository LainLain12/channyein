package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "channyein.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// Optionally, you can ping the database to ensure it's working
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}
