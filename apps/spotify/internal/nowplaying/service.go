package nowplaying

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
)

const defaultSpotifyAPIBaseURL = "https://api.spotify.com/v1"

type Service interface {
	GetCurrent(ctx context.Context) (Response, error)
}

type SpotifyService struct {
	client  *http.Client
	baseURL string
}

func NewSpotifyService(client *http.Client, baseURL string) *SpotifyService {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	if baseURL == "" {
		baseURL = defaultSpotifyAPIBaseURL
	}
	return &SpotifyService{
		client:  client,
		baseURL: baseURL,
	}
}

func (s *SpotifyService) GetCurrent(ctx context.Context) (Response, error) {
	token, ok := auth.AccessTokenFromContext(ctx)
	if !ok || token == "" {
		return Response{}, auth.ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/me/player/currently-playing", nil)
	if err != nil {
		return Response{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("spotify request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return Response{IsPlaying: false}, nil
	case http.StatusOK:
		var payload struct {
			IsPlaying  bool `json:"is_playing"`
			ProgressMs int  `json:"progress_ms"`
			Item       struct {
				Name       string `json:"name"`
				DurationMs int    `json:"duration_ms"`
				Artists    []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Images []struct {
						URL string `json:"url"`
					} `json:"images"`
				} `json:"album"`
			} `json:"item"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return Response{}, fmt.Errorf("decode response: %w", err)
		}

		artists := make([]string, 0, len(payload.Item.Artists))
		for _, artist := range payload.Item.Artists {
			if artist.Name != "" {
				artists = append(artists, artist.Name)
			}
		}

		albumCoverURL := ""
		if len(payload.Item.Album.Images) > 0 {
			albumCoverURL = payload.Item.Album.Images[0].URL
		}

		return Response{
			IsPlaying:     payload.IsPlaying,
			Track:         payload.Item.Name,
			Artists:       artists,
			AlbumCoverURL: albumCoverURL,
			ProgressMs:    payload.ProgressMs,
			DurationMs:    payload.Item.DurationMs,
		}, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return Response{}, auth.ErrUnauthorized
	default:
		return Response{}, fmt.Errorf("spotify auth failed: status %d", resp.StatusCode)
	}
}
