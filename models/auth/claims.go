package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}
