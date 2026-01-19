package auth

type TokenRequest struct {
	Code     string
	Provider string
}
