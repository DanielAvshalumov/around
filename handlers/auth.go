package handlers

import (

	// "fmt"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielavshalumov/around/models/auth"
	"github.com/danielavshalumov/around/services"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	UserService *services.UserService
}

func NewAuthHandler(us *services.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: us,
	}
}

func (a *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie == nil {
		http.Error(w, "Error getting Session Cookie"+err.Error(), 401)
		return
	}
	cookieString := sessionCookie.Value
	claims, err := validateJwt(cookieString)
	if err != nil {
		http.Error(w, "Error validating JWT", 401)
		return
	}
	fmt.Printf("Claims %+v\n", claims)
	user, err := a.UserService.GetUserById(claims.UserId)
	if err != nil {
		http.Error(w, "User not found", 500)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (a *AuthHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	json.NewEncoder(w).Encode(models.SimpleError{
	// 		Error: "Method Not Allowed",
	// 	})
	// 	return
	// }
	code := r.URL.Query().Get("code")
	provider := "google"
	tokenReq := auth.TokenRequest{Code: code, Provider: provider}
	// err := json.NewDecoder(r.Body).Decode(&tokenReq)

	if tokenReq.Code == "" {
		http.Error(w, "missing auth code", http.StatusBadRequest)
	}
	tokenRes, err := services.GetToken(tokenReq.Code)
	if err != nil {
		http.Error(w, "failed making Google POST request: "+err.Error(), 500)
	}
	userDTO, err := services.GetUserInfo(tokenRes)
	if err != nil {
		http.Error(w, "failed getting user Info"+err.Error(), 500)
	}

	user, err := a.UserService.GetUserFromOAuth(*userDTO)

	sessionToken, err := generateJwt(user.UserId, user.Email, tokenReq.Provider)
	if err != nil {
		http.Error(w, "failed JWT step"+err.Error(), 500)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800,
		Path:     "/",
	})

	http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
}

func generateJwt(userId int64, email string, provider string) (string, error) {
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

func validateJwt(tokenString string) (*auth.Claims, error) {
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
