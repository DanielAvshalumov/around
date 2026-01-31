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

func (db *Db) GetBacklinkById(backlinkId int) (*models.Backlink, error) {
	backlink := &models.Backlink{}
	query := `SELECT id, source, link, s_id, title from backlinks where id = ?`
	err := db.QueryRow(query, backlinkId).Scan(&backlink.Id, &backlink.Source, &backlink.Link, &backlink.Session, &backlink.Title)
	if err != nil {
		fmt.Println("Error in config/db", err.Error())
		return nil, err
	}
	return backlink, nil
}

func (db *Db) InsertIntoBacklink(backlink *models.Backlink) (int64, error) {
	query := `
		INSERT INTO backlinks (source, link, dofollow, s_id, title)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := db.Exec(query, backlink.Source, backlink.Link, backlink.Dofollow, backlink.Session, backlink.Title)
	if err != nil {
		fmt.Printf("DB Error - Failed to insert backlinks %s -> %s", backlink.Source, backlink.Link)
		return 0, err
	}
	return res.LastInsertId()
}

func (db *Db) InsertUserBacklink(backlink_id int, user_id int64, response string) (int64, error) {
	query := `
		INSERT INTO user_backlink (user_id, backlink_id, response)
		VALUES (?, ?, ?)
	`
	res, err := db.Exec(query, user_id, backlink_id, response)
	if err != nil {
		fmt.Println("DB Error - failed to insert user_backlink", err.Error())
	}
	return res.RowsAffected()
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
	query := `SELECT userId, googleSub, email, dateCreated from users where userId = ?`
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

func (db *Db) GetUserBacklinks(userId int64) ([]models.Backlink, error) {
	var userBacklinks []models.Backlink
	query := `SELECT b.* FROM backlinks b inner join user_backlink ub on ub.backlink_id = b.id where ub.user_id = ?`

	rows, err := db.Query(query, userId)
	if err != nil {
		fmt.Println("failed to get rows from user_backlink")
	}
	defer rows.Close()

	for rows.Next() {
		var userBacklink models.Backlink
		err = rows.Scan(
			&userBacklink.Id,
			&userBacklink.Source,
			&userBacklink.Link,
			&userBacklink.Dofollow,
			&userBacklink.Session,
			&userBacklink.Title,
		)
		if err != nil {
			fmt.Println("Error seriazliing a row for User backlinks")
			return nil, err
		}
		userBacklinks = append(userBacklinks, userBacklink)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Error after looping through all rows for user backlinks")
	}
	return userBacklinks, nil
}
