package handlers

import (

	// "fmt"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/lib"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/models/auth"
	"github.com/danielavshalumov/around/services"
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
	claims, err := lib.ValidateJwt(cookieString)
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
	fmt.Println("user", user)
	json.NewEncoder(w).Encode(user)
}

func (a *AuthHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method Not Allowed",
		})
		return
	}
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

	sessionToken, err := lib.GenerateJwt(user.UserId, user.Email, tokenReq.Provider)
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
