package lib

import (
	"fmt"
	"net/http"
)

func GetUserIdFromCookie(r *http.Request) (int64, error) {
	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie == nil {
		return 0, err
	}
	cookieString := sessionCookie.Value
	claims, err := ValidateJwt(cookieString)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Claims %+v\n", claims)
	return claims.UserId, nil
}
