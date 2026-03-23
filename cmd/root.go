package cmd

import (
	"os"

	"github.com/jrogala/spotify-cli/config"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "spotify",
	Short: "A command-line interface for controlling Spotify playback",
	Long: `A command-line interface for controlling Spotify playback.

Quick examples:
  spotify play --uri spotify:album:XXX --device Cuisine   Play an album on a device
  spotify play --uri spotify:artist:XXX --device Cuisine   Play an artist radio
  spotify play                                             Resume playback
  spotify pause                                            Pause playback
  spotify next / prev                                      Skip tracks
  spotify now-playing                                      Show current track
  spotify search "query" --type album                      Search (track,album,artist)
  spotify volume 30                                        Set volume (0-100)
  spotify devices                                          List Spotify Connect devices
  spotify transfer --device "device name"                  Move playback to a device
  spotify queue                                            Show playback queue
  spotify top --type tracks                                Your top tracks
  spotify top --type artists                               Your top artists`,
}

func Execute() {
	config.Init()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output raw JSON responses")
}
