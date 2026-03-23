package ops

import (
	"fmt"

	"github.com/jrogala/spotify-cli/client"
)

// VolumeInfo holds volume state.
type VolumeInfo struct {
	Level  int    `json:"volume"`
	Device string `json:"device"`
}

// GetVolume returns the current volume.
func GetVolume(c *client.Client) (*VolumeInfo, error) {
	state, err := c.GetPlaybackState()
	if err != nil {
		return nil, err
	}
	return &VolumeInfo{
		Level:  state.Device.VolumePercent,
		Device: state.Device.Name,
	}, nil
}

// SetVolume sets the volume (0-100).
func SetVolume(c *client.Client, deviceID string, level int) error {
	if level < 0 || level > 100 {
		return fmt.Errorf("volume must be 0-100")
	}
	return c.SetVolume(deviceID, level)
}
