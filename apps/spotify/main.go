package main

import (
	"log"
	"net/http"
	"os"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/api"
	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
	"github.com/harrytjbreen/homelab/apps/spotify/internal/nowplaying"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	redirectURL := os.Getenv("SPOTIFY_REDIRECT_URL")
	cookieSecret := os.Getenv("SPOTIFY_SESSION_SECRET")
	if clientID == "" || clientSecret == "" || redirectURL == "" || cookieSecret == "" {
		log.Printf("warning: missing Spotify OAuth settings")
	}

	oauth := auth.NewSpotifyOAuth(nil, clientID, clientSecret, redirectURL)
	cookieCodec, err := auth.NewCookieCodec(cookieSecret)
	if err != nil {
		log.Fatalf("failed to configure session cookies: %v", err)
	}

	mux := http.NewServeMux()

	nowPlayingService := nowplaying.NewSpotifyService(nil, "")
	handler := api.NewHandler(nowPlayingService, oauth, cookieCodec)
	handler.Register(mux)

	addr := os.Getenv("SPOTIFY_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
