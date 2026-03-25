package api

import (
	"net/http"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
	"github.com/harrytjbreen/homelab/apps/spotify/internal/nowplaying"
)

type Handler struct {
	nowPlayingService nowplaying.Service
	oauth             *auth.SpotifyOAuth
	cookieCodec       *auth.CookieCodec
}

func NewHandler(nowPlayingService nowplaying.Service, oauth *auth.SpotifyOAuth, cookieCodec *auth.CookieCodec) *Handler {
	return &Handler{
		nowPlayingService: nowPlayingService,
		oauth:             oauth,
		cookieCodec:       cookieCodec,
	}
}

const basePath = "/spotify"

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc(basePath+"/health", h.health)
	mux.HandleFunc(basePath+"/auth/login", h.login)
	mux.HandleFunc(basePath+"/auth/callback", h.callback)
	mux.Handle(basePath+"/api/now-playing", h.requireAuth(http.HandlerFunc(h.nowPlaying)))
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
