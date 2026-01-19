package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/danielavshalumov/around/models/auth"
)

type AuthService struct {
}

func GetToken(code string) (*auth.TokenResponse, error) {
	// Verify User
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("OAUTH_REDIRECT_URI"))
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token?grant_type=authorization_code", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var token auth.TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GetUserInfo(token *auth.TokenResponse) (*auth.UserDTO, error) {
	var user *auth.UserDTO
	fmt.Println("Token", token)
	res, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token.IDToken))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &user)
	fmt.Println(res.StatusCode)
	if err != nil {
		return nil, err
	}
	fmt.Println("user", user)
	return user, nil
}
