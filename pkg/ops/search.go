package ops

import (
	"github.com/jrogala/spotify-cli/client"
)

// TrackResult holds a search track result.
type TrackResult struct {
	URI      string `json:"uri"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration string `json:"duration"`
}

// AlbumResult holds a search album result.
type AlbumResult struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	ReleaseDate string `json:"release_date"`
	TotalTracks int    `json:"total_tracks"`
}

// ArtistResult holds a search artist result.
type ArtistResult struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// SearchResults holds all search results.
type SearchResults struct {
	Tracks  []TrackResult  `json:"tracks,omitempty"`
	Albums  []AlbumResult  `json:"albums,omitempty"`
	Artists []ArtistResult `json:"artists,omitempty"`
}

// Search searches Spotify.
func Search(c *client.Client, query, types string, limit int) (*SearchResults, error) {
	result, err := c.Search(query, types, limit)
	if err != nil {
		return nil, err
	}

	sr := &SearchResults{}

	if result.Tracks != nil {
		sr.Tracks = make([]TrackResult, len(result.Tracks.Items))
		for i, t := range result.Tracks.Items {
			sr.Tracks[i] = TrackResult{
				URI:      t.URI,
				Name:     t.Name,
				Artist:   client.ArtistNames(t.Artists),
				Album:    t.Album.Name,
				Duration: client.FormatDuration(t.DurationMs),
			}
		}
	}

	if result.Albums != nil {
		sr.Albums = make([]AlbumResult, len(result.Albums.Items))
		for i, a := range result.Albums.Items {
			sr.Albums[i] = AlbumResult{
				URI:         a.URI,
				Name:        a.Name,
				Artist:      client.ArtistNames(a.Artists),
				ReleaseDate: a.ReleaseDate,
				TotalTracks: a.TotalTracks,
			}
		}
	}

	if result.Artists != nil {
		sr.Artists = make([]ArtistResult, len(result.Artists.Items))
		for i, a := range result.Artists.Items {
			sr.Artists[i] = ArtistResult{URI: a.URI, Name: a.Name}
		}
	}

	return sr, nil
}
