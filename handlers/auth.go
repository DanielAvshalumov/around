package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
)

func HandleAuthStart(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Auth Error %v", err)
		http.Redirect(w, r, fmt.Sprintf("http://localhost:3000/auth/error?messages=%s", err.Error()), http.StatusTemporaryRedirect)
		return
	}

}

func generateJwt(userId, email, name, provider string) (string, error) {
	claims := Claims{
		UserId:   userId,
		Email:    email,
		Name:     name,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}
