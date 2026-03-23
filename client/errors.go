package client

import (
	"encoding/json"
	"fmt"
	"io"
)

// SpotifyError represents an error returned by the Spotify API.
type SpotifyError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *SpotifyError) Error() string {
	return fmt.Sprintf("spotify error %d: %s", e.Status, e.Message)
}

func parseError(body io.Reader) error {
	var envelope struct {
		Error SpotifyError `json:"error"`
	}
	if err := json.NewDecoder(body).Decode(&envelope); err != nil {
		return fmt.Errorf("unknown spotify error")
	}
	return &envelope.Error
}
