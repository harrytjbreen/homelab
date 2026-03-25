package api

import (
	"net/http"
	"time"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
)

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.oauth == nil || h.cookieCodec == nil {
			unauthorized(w)
			return
		}

		authCookie, err := r.Cookie(sessionCookieName)
		if err != nil || authCookie.Value == "" {
			unauthorized(w)
			return
		}

		session, err := h.cookieCodec.Decode(authCookie.Value)
		if err != nil {
			unauthorized(w)
			return
		}

		if needsRefresh(session) {
			token, err := h.oauth.RefreshToken(r.Context(), session.RefreshToken)
			if err != nil {
				unauthorized(w)
				return
			}

			session.AccessToken = token.AccessToken
			if token.RefreshToken != "" {
				session.RefreshToken = token.RefreshToken
			}
			if !token.ExpiresAt.IsZero() {
				session.ExpiresAt = token.ExpiresAt
			}

			cookieValue, err := h.cookieCodec.Encode(session)
			if err != nil {
				unauthorized(w)
				return
			}
			http.SetCookie(w, sessionCookie(r, cookieValue, session.ExpiresAt))
		}

		ctx := auth.WithUser(r.Context(), session.User)
		ctx = auth.WithAccessToken(ctx, session.AccessToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func needsRefresh(session auth.SessionData) bool {
	if session.AccessToken == "" {
		return true
	}
	if session.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(session.ExpiresAt.Add(-30 * time.Second))
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
