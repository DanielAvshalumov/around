package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	displayOpeningMessage()

	db, err := sql.Open("sqlite3", "./avsolutions.db")
	if err != nil {
		return nil, fmt.Errorf("failed opening database: %w", err)
	}

	createBacklinkTable(db)

	return db, nil
}

func createBacklinkTable(db *sql.DB) {
	fmt.Println("Opening backlinks table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS backlinks (
		id INTEGER PRIMARY KEY,
		source TEXT NOT NULL,
		link TEXT NOT NULL,
		dofollow INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		fmt.Println("Failed to create table:", err)
	}
}

func displayOpeningMessage() {
	fmt.Println("Opening Sqlite3")
	for i := 0; i < 6; i++ {
		dots := i
		fmt.Printf("\rLoading%s", strings.Repeat(".", dots))
		time.Sleep(500 * time.Millisecond)
	}
}
