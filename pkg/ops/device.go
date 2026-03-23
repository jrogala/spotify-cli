// Package ops contains business logic for Spotify operations.
// Functions return Go structs and errors — zero I/O, zero formatting.
package ops

import (
	"fmt"
	"strings"

	"github.com/jrogala/spotify-cli/client"
)

// DeviceInfo holds device details.
type DeviceInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
	Volume   int    `json:"volume_percent"`
}

// ListDevices returns all Spotify Connect devices.
func ListDevices(c *client.Client) ([]DeviceInfo, error) {
	devices, err := c.GetDevices()
	if err != nil {
		return nil, err
	}
	result := make([]DeviceInfo, len(devices))
	for i, d := range devices {
		result[i] = DeviceInfo{
			ID:       d.ID,
			Name:     d.Name,
			Type:     d.Type,
			IsActive: d.IsActive,
			Volume:   d.VolumePercent,
		}
	}
	return result, nil
}

// ResolveDeviceID resolves a device name to its ID. Empty string returns empty.
func ResolveDeviceID(c *client.Client, nameOrID string) (string, error) {
	if nameOrID == "" {
		return "", nil
	}
	if len(nameOrID) > 30 {
		return nameOrID, nil
	}
	devices, err := c.GetDevices()
	if err != nil {
		return "", fmt.Errorf("listing devices to resolve name: %w", err)
	}
	lower := strings.ToLower(nameOrID)
	for _, d := range devices {
		if strings.ToLower(d.Name) == lower {
			return d.ID, nil
		}
	}
	return "", fmt.Errorf("device %q not found", nameOrID)
}
