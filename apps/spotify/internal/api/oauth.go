package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
)

const (
	sessionCookieName = "spotify_session"
	stateCookieName   = "spotify_oauth_state"
)

var spotifyScopes = []string{
	"user-read-currently-playing",
	"user-read-playback-state",
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	state, err := randomState(24)
	if err != nil {
		http.Error(w, "failed to start auth", http.StatusInternalServerError)
		return
	}

	authURL, err := h.oauth.AuthURL(state, spotifyScopes)
	if err != nil {
		http.Error(w, "failed to build auth url", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     basePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Error(w, "missing state or code", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil || stateCookie.Value != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	token, err := h.oauth.ExchangeCode(r.Context(), code)
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusUnauthorized)
		return
	}

	user, err := h.oauth.GetUser(r.Context(), token.AccessToken)
	if err != nil {
		http.Error(w, "failed to validate user", http.StatusUnauthorized)
		return
	}

	session := auth.SessionData{
		User:         user,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	cookieValue, err := h.cookieCodec.Encode(session)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, sessionCookie(r, cookieValue, session.ExpiresAt))

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     basePath,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	if redirectURL := os.Getenv("FRONTEND_REDIRECT_URL"); redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("authorized"))
}

func sessionCookie(r *http.Request, value string, expiresAt time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     basePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}
	if r.TLS != nil {
		cookie.Secure = true
	}
	return cookie
}

func randomState(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
