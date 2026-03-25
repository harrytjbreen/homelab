package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
)

var ErrInvalidCookie = errors.New("invalid session cookie")

type CookieCodec struct {
	secret []byte
}

func NewCookieCodec(secret string) (*CookieCodec, error) {
	if len(secret) < 32 {
		return nil, errors.New("cookie secret must be at least 32 characters")
	}
	return &CookieCodec{secret: []byte(secret)}, nil
}

func (c *CookieCodec) Encode(session SessionData) (string, error) {
	payload, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	sig := c.sign(payload)
	return base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func (c *CookieCodec) Decode(value string) (SessionData, error) {
	parts := splitOnce(value, '.')
	if len(parts) != 2 {
		return SessionData{}, ErrInvalidCookie
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return SessionData{}, ErrInvalidCookie
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return SessionData{}, ErrInvalidCookie
	}

	if !hmac.Equal(sig, c.sign(payload)) {
		return SessionData{}, ErrInvalidCookie
	}

	var session SessionData
	if err := json.Unmarshal(payload, &session); err != nil {
		return SessionData{}, ErrInvalidCookie
	}

	return session, nil
}

func (c *CookieCodec) sign(payload []byte) []byte {
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write(payload)
	return mac.Sum(nil)
}

func splitOnce(value string, sep byte) []string {
	for i := 0; i < len(value); i++ {
		if value[i] == sep {
			return []string{value[:i], value[i+1:]}
		}
	}
	return nil
}
