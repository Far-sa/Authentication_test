package auth

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tokens.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create tokens table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create tokens table: %v", err)
	}
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}
