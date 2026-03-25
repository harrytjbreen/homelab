package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/harrytjbreen/homelab/apps/spotify/internal/auth"
)

func (h *Handler) nowPlaying(w http.ResponseWriter, r *http.Request) {
	response, err := h.nowPlayingService.GetCurrent(r.Context())
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "failed to get now playing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
