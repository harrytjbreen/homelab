package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	DefaultAuthURL    = "https://accounts.spotify.com/authorize"
	DefaultTokenURL   = "https://accounts.spotify.com/api/token"
	DefaultAPIBaseURL = "https://api.spotify.com/v1"
)

type SpotifyOAuth struct {
	client     *http.Client
	config     *oauth2.Config
	apiBaseURL string
}

func NewSpotifyOAuth(client *http.Client, clientID, clientSecret, redirectURL string) *SpotifyOAuth {
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   DefaultAuthURL,
		TokenURL:  DefaultTokenURL,
		AuthStyle: oauth2.AuthStyleInHeader,
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     endpoint,
	}

	return &SpotifyOAuth{
		client:     client,
		config:     config,
		apiBaseURL: DefaultAPIBaseURL,
	}
}

func (s *SpotifyOAuth) AuthURL(state string, scopes []string) (string, error) {
	if s.config.ClientID == "" || s.config.RedirectURL == "" {
		return "", fmt.Errorf("missing client id or redirect url")
	}

	options := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	if len(scopes) > 0 {
		options = append(options, oauth2.SetAuthURLParam("scope", strings.Join(scopes, " ")))
	}

	return s.config.AuthCodeURL(state, options...), nil
}

func (s *SpotifyOAuth) ExchangeCode(ctx context.Context, code string) (Token, error) {
	if s.config.ClientID == "" || s.config.ClientSecret == "" || s.config.RedirectURL == "" {
		return Token{}, fmt.Errorf("missing client credentials or redirect url")
	}

	oauthToken, err := s.config.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return Token{}, fmt.Errorf("token request failed: %w", err)
	}

	if oauthToken.AccessToken == "" {
		return Token{}, ErrUnauthorized
	}

	return Token{
		AccessToken:  oauthToken.AccessToken,
		TokenType:    oauthToken.TokenType,
		RefreshToken: oauthToken.RefreshToken,
		ExpiresAt:    oauthToken.Expiry,
	}, nil
}

func (s *SpotifyOAuth) GetUser(ctx context.Context, accessToken string) (User, error) {
	if accessToken == "" {
		return User{}, ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiBaseURL+"/me", nil)
	if err != nil {
		return User{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("spotify request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var payload User
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return User{}, fmt.Errorf("decode response: %w", err)
		}
		if payload.ID == "" {
			return User{}, ErrUnauthorized
		}
		return payload, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return User{}, ErrUnauthorized
	default:
		return User{}, fmt.Errorf("spotify auth failed: status %d", resp.StatusCode)
	}
}

func (s *SpotifyOAuth) RefreshToken(ctx context.Context, refreshToken string) (Token, error) {
	if s.config.ClientID == "" || s.config.ClientSecret == "" {
		return Token{}, fmt.Errorf("missing client credentials")
	}
	if refreshToken == "" {
		return Token{}, ErrUnauthorized
	}

	source := s.config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-1 * time.Hour),
	})

	oauthToken, err := source.Token()
	if err != nil {
		return Token{}, fmt.Errorf("refresh token failed: %w", err)
	}
	if oauthToken.AccessToken == "" {
		return Token{}, ErrUnauthorized
	}

	return Token{
		AccessToken:  oauthToken.AccessToken,
		TokenType:    oauthToken.TokenType,
		RefreshToken: oauthToken.RefreshToken,
		ExpiresAt:    oauthToken.Expiry,
	}, nil
}
