package config

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/danielavshalumov/around/models"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	client *sql.DB
}

func NewDB(db *sql.DB) *Db {
	return &Db{
		client: db,
	}
}

func InitDB() (*Db, error) {
	displayOpeningMessage()

	db, err := sql.Open("sqlite3", "./config/avsolutions.db")
	if err != nil {
		return nil, fmt.Errorf("failed opening database: %w", err)
	}
	Db := NewDB(db)
	createBacklinkTable(Db.client)

	return Db, nil
}

func (db *Db) InsertIntoBacklink(backlink *models.Backlink) (int64, error) {
	query := `
		INSERT INTO backlinks (source, link, dofollow)
		VALUES (?, ?, ?)
	`
	res, err := db.client.Exec(query, backlink.Source, backlink.Link, backlink.Dofollow)
	if err != nil {
		fmt.Printf("DB Error - Failed to insert backlinks %s -> %s", backlink.Source, backlink.Link)
		return 0, err
	}
	return res.LastInsertId()
}

func createBacklinkTable(db *sql.DB) {
	fmt.Println("Opening backlinks table")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS backlinks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source TEXT NOT NULL,
		link TEXT NOT NULL,
		dofollow INTEGER NOT NULL
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
