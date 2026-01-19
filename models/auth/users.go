package auth

import (
	"time"
)

type User struct {
	UserId     int64
	GoogleSub  string
	Email      string
	DateCreate time.Time
}

type UserDTO struct {
	Sub   string `json:"sub"`
	Email string
}
