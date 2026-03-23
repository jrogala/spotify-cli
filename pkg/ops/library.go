package ops

import (
	"github.com/jrogala/spotify-cli/client"
)

// TopTrack holds a top track entry.
type TopTrack struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	URI    string `json:"uri"`
}

// TopArtist holds a top artist entry.
type TopArtist struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// LikedTrack holds a saved track entry.
type LikedTrack struct {
	Name    string `json:"name"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	AddedAt string `json:"added_at"`
}

// LikedResult holds liked tracks with pagination info.
type LikedResult struct {
	Total int          `json:"total"`
	Items []LikedTrack `json:"items"`
}

// RecentTrack holds a recently played track.
type RecentTrack struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	PlayedAt string `json:"played_at"`
}

// GetTopTracks returns the user's top tracks.
func GetTopTracks(c *client.Client, timeRange string, limit int) ([]TopTrack, error) {
	tracks, err := c.GetTopTracks(timeRange, limit)
	if err != nil {
		return nil, err
	}
	result := make([]TopTrack, len(tracks))
	for i, t := range tracks {
		result[i] = TopTrack{
			Name:   t.Name,
			Artist: client.ArtistNames(t.Artists),
			Album:  t.Album.Name,
			URI:    t.URI,
		}
	}
	return result, nil
}

// GetTopArtists returns the user's top artists.
func GetTopArtists(c *client.Client, timeRange string, limit int) ([]TopArtist, error) {
	artists, err := c.GetTopArtists(timeRange, limit)
	if err != nil {
		return nil, err
	}
	result := make([]TopArtist, len(artists))
	for i, a := range artists {
		result[i] = TopArtist{Name: a.Name, URI: a.URI}
	}
	return result, nil
}

// GetLikedTracks returns the user's saved tracks.
func GetLikedTracks(c *client.Client, limit, offset int) (*LikedResult, error) {
	paging, err := c.GetLikedTracks(limit, offset)
	if err != nil {
		return nil, err
	}
	items := make([]LikedTrack, len(paging.Items))
	for i, st := range paging.Items {
		items[i] = LikedTrack{
			Name:    st.Track.Name,
			Artist:  client.ArtistNames(st.Track.Artists),
			Album:   st.Track.Album.Name,
			AddedAt: st.AddedAt,
		}
	}
	return &LikedResult{Total: paging.Total, Items: items}, nil
}

// GetRecentlyPlayed returns recently played tracks.
func GetRecentlyPlayed(c *client.Client, limit int) ([]RecentTrack, error) {
	history, err := c.GetRecentlyPlayed(limit)
	if err != nil {
		return nil, err
	}
	result := make([]RecentTrack, len(history))
	for i, h := range history {
		result[i] = RecentTrack{
			Name:     h.Track.Name,
			Artist:   client.ArtistNames(h.Track.Artists),
			PlayedAt: h.PlayedAt,
		}
	}
	return result, nil
}
