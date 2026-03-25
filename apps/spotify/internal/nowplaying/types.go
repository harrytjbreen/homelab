package nowplaying

type Response struct {
	IsPlaying     bool     `json:"is_playing"`
	Track         string   `json:"track"`
	Artists       []string `json:"artists"`
	AlbumCoverURL string   `json:"album_cover_url"`
	ProgressMs    int      `json:"progress_ms"`
	DurationMs    int      `json:"duration_ms"`
}
