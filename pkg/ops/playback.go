package ops

import (
	"github.com/jrogala/spotify-cli/client"
)

// NowPlaying holds the current playback state.
type NowPlaying struct {
	IsPlaying bool   `json:"is_playing"`
	Track     string `json:"track"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Progress  string `json:"progress"`
	Duration  string `json:"duration"`
	Device    string `json:"device"`
	Volume    int    `json:"volume"`
}

// GetNowPlaying returns the current playback state.
func GetNowPlaying(c *client.Client) (*NowPlaying, error) {
	state, err := c.GetPlaybackState()
	if err != nil {
		return nil, err
	}
	return &NowPlaying{
		IsPlaying: state.IsPlaying,
		Track:     state.Item.Name,
		Artist:    client.ArtistNames(state.Item.Artists),
		Album:     state.Item.Album.Name,
		Progress:  client.FormatDuration(state.ProgressMs),
		Duration:  client.FormatDuration(state.Item.DurationMs),
		Device:    state.Device.Name,
		Volume:    state.Device.VolumePercent,
	}, nil
}

// Play starts or resumes playback.
func Play(c *client.Client, deviceID, contextURI string, uris []string) error {
	return c.Play(deviceID, contextURI, uris, -1)
}

// Pause pauses playback.
func Pause(c *client.Client, deviceID string) error {
	return c.Pause(deviceID)
}

// Next skips to the next track.
func Next(c *client.Client, deviceID string) error {
	return c.Next(deviceID)
}

// Previous goes back to the previous track.
func Previous(c *client.Client, deviceID string) error {
	return c.Previous(deviceID)
}

// TransferPlayback moves playback to a device.
func TransferPlayback(c *client.Client, deviceID string, play bool) error {
	return c.TransferPlayback(deviceID, play)
}
