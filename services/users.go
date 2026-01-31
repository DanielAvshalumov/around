package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/models/auth"
)

type UserService struct {
	DB *config.Db
}

func CreateUserService(db *config.Db) *UserService {
	return &UserService{
		DB: db,
	}
}

func (u UserService) GetUserById(userId int64) (*auth.User, error) {
	fmt.Println(userId)
	user, err := u.DB.GetUserById(userId)
	if err != nil {
		fmt.Println("Error in User service getting User from UserId", err)
		return nil, err
	}
	return user, nil
}

func (u UserService) GetUserFromOAuth(dto auth.UserDTO) (*auth.User, error) {
	user, err := u.DB.GetUser(dto.Sub)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Database Error", err)
	}
	fmt.Println("existing user", user)
	if user != nil {
		return user, nil
	}
	newUser, err := u.DB.CreateUser(dto.Email, dto.Sub)
	newUser.Email = dto.Email
	newUser.GoogleSub = dto.Sub
	if err != nil {
		fmt.Println("Error Creating User", err)
	}
	return newUser, nil
}

func (u UserService) SaveBacklink(backlink_id int, user_id int64, response string) (int64, error) {
	id, err := u.DB.InsertUserBacklink(backlink_id, user_id, response)
	if err != nil {
		fmt.Println("Error DB inserting", err.Error())
		return 0, err
	}
	return id, nil
}

func (u UserService) GetUserBacklinks(userId int64) ([]models.Backlink, error) {
	userBacklinks, err := u.DB.GetUserBacklinks(userId)
	if err != nil {
		fmt.Println("UserService - error getting user bckloinks from db call")
		return nil, err
	}
	return userBacklinks, nil
}
