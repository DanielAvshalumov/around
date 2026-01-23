package lib

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/danielavshalumov/around/models/auth"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(userId int64, email string, provider string) (string, error) {
	claims := auth.Claims{
		UserId:   userId,
		Email:    email,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJwt(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println("Failed parsing claims")
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return nil, errors.New("Invalid claims")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Expired token")
	}

	return claims, nil
}
