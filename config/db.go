package config

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/models/auth"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	*sql.DB
}

func NewDB(db *sql.DB) *Db {
	return &Db{
		db,
	}
}

func InitDB() (*Db, error) {
	displayOpeningMessage()

	db, err := sql.Open("sqlite3", "./config/avsolutions.db?_loc=auto&parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed opening database: %w", err)
	}
	Db := NewDB(db)
	createBacklinkTable(Db.DB)

	return Db, nil
}

func (db *Db) CreateNewCrawlSession(userId int64) (int64, error) {
	query := "INSERT INTO crawl_session (user_id) VALUES (?)"
	res, err := db.Exec(query, userId)
	if err != nil {
		fmt.Println("DB Error - Failed to create session", err.Error())
		return 0, err
	}
	return res.LastInsertId()
}

func (db *Db) InsertIntoBacklink(backlink *models.Backlink) (int64, error) {
	query := `
		INSERT INTO backlinks (source, link, dofollow, s_id)
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, backlink.Source, backlink.Link, backlink.Dofollow, backlink.Session)
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

func (db *Db) GetUserById(userId int64) (*auth.User, error) {
	user := &auth.User{}
	var sqlDate string
	query := `SELECT userId, googleSub, email, dateCreated from users where userId`
	err := db.QueryRow(query, userId).Scan(&user.UserId, &user.GoogleSub, &user.Email, &sqlDate)
	if err != nil {
		fmt.Println("User not found by Id")
		fmt.Println(err)
		return nil, err
	}
	convertedDate, err := time.Parse("2006-01-02 15:04:05", sqlDate)
	user.DateCreate = convertedDate
	return user, nil
}

func (db *Db) GetUser(googleSub string) (*auth.User, error) {
	user := &auth.User{}
	var sqlDate string
	query := `SELECT userId, googleSub, email, dateCreated from users where googleSub = ?`
	err := db.QueryRow(query, googleSub).Scan(&user.UserId, &user.GoogleSub, &user.Email, &sqlDate)
	convertedDate, err := time.Parse("2006-01-02 15:04:05", sqlDate)
	user.DateCreate = convertedDate
	if err != nil {
		fmt.Println("User not found", err)
		return nil, err
	}
	return user, nil
}

func (db *Db) CreateUser(email string, googleSub string) (*auth.User, error) {
	user := &auth.User{}
	query := `INSERT INTO users (googleSub, email) VALUES (?, ?)`
	result, err := db.Exec(query, googleSub, email)
	if err != nil {
		fmt.Println("User already exists")
		return nil, err
	}
	id, err := result.LastInsertId()
	user.UserId = id
	return user, nil
}
