package client

// Device represents a Spotify Connect device.
type Device struct {
	ID            string `json:"id"`
	IsActive      bool   `json:"is_active"`
	IsRestricted  bool   `json:"is_restricted"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	VolumePercent int    `json:"volume_percent"`
	SupportsVolume bool  `json:"supports_volume"`
}

// Artist represents a Spotify artist.
type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// Image represents an image resource.
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Album represents a Spotify album.
type Album struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Artists     []Artist `json:"artists"`
	Images      []Image  `json:"images"`
	URI         string   `json:"uri"`
	ReleaseDate string   `json:"release_date"`
	TotalTracks int      `json:"total_tracks"`
}

// Track represents a Spotify track.
type Track struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []Artist `json:"artists"`
	Album      Album    `json:"album"`
	DurationMs int     `json:"duration_ms"`
	URI        string   `json:"uri"`
	TrackNumber int    `json:"track_number"`
}

// PlaybackState represents the current playback state.
type PlaybackState struct {
	Device       Device `json:"device"`
	IsPlaying    bool   `json:"is_playing"`
	Item         Track  `json:"item"`
	ProgressMs   int    `json:"progress_ms"`
	ShuffleState bool   `json:"shuffle_state"`
	RepeatState  string `json:"repeat_state"`
}

// QueueResponse represents the user's playback queue.
type QueueResponse struct {
	CurrentlyPlaying Track   `json:"currently_playing"`
	Queue            []Track `json:"queue"`
}

// Paging represents a paginated response.
type Paging[T any] struct {
	Items    []T    `json:"items"`
	Total    int    `json:"total"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// SavedTrack wraps a track with the date it was saved.
type SavedTrack struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

// PlayHistory represents a recently played item.
type PlayHistory struct {
	Track    Track  `json:"track"`
	PlayedAt string `json:"played_at"`
}

// SearchResult represents search results.
type SearchResult struct {
	Tracks  *Paging[Track]  `json:"tracks,omitempty"`
	Artists *Paging[Artist] `json:"artists,omitempty"`
	Albums  *Paging[Album]  `json:"albums,omitempty"`
}
