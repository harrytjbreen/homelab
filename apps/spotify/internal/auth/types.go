package auth

import (
	"errors"
	"time"
)

var ErrUnauthorized = errors.New("unauthorized")

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresAt    time.Time
}

type SessionData struct {
	User         User      `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
