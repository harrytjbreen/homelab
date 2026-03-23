package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type pingResponse struct {
	Message string `json:"message"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pingResponse{Message: "pong"})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)

	addr := os.Getenv("PING_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("ping api listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
