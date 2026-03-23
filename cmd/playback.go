package cmd

import (
	"fmt"

	"github.com/jrogala/spotify-cli/config"
	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playCmd, pauseCmd, nextCmd, prevCmd, nowPlayingCmd, transferCmd)

	playCmd.Flags().String("device", "", "Target device ID or name")
	playCmd.Flags().String("uri", "", "Spotify URI (album, playlist, or track)")
	playCmd.Flags().StringSlice("tracks", nil, "Track URIs to play")

	pauseCmd.Flags().String("device", "", "Target device ID or name")
	nextCmd.Flags().String("device", "", "Target device ID or name")
	prevCmd.Flags().String("device", "", "Target device ID or name")

	transferCmd.Flags().String("device", "", "Target device ID or name")
	transferCmd.MarkFlagRequired("device")
	transferCmd.Flags().Bool("play", true, "Start playback on transfer")
}

func resolveDevice(cmd *cobra.Command) string {
	dev, _ := cmd.Flags().GetString("device")
	if dev != "" {
		return dev
	}
	return config.DefaultDevice()
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start or resume playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		uri, _ := cmd.Flags().GetString("uri")
		tracks, _ := cmd.Flags().GetStringSlice("tracks")

		if err := ops.Play(c, deviceID, uri, tracks); err != nil {
			return err
		}
		fmt.Println("Playing")
		return nil
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		if err := ops.Pause(c, deviceID); err != nil {
			return err
		}
		fmt.Println("Paused")
		return nil
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to next track",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		if err := ops.Next(c, deviceID); err != nil {
			return err
		}
		fmt.Println("Next")
		return nil
	},
}

var prevCmd = &cobra.Command{
	Use:     "prev",
	Aliases: []string{"previous"},
	Short:   "Skip to previous track",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		if err := ops.Previous(c, deviceID); err != nil {
			return err
		}
		fmt.Println("Previous")
		return nil
	},
}

var nowPlayingCmd = &cobra.Command{
	Use:     "now-playing",
	Aliases: []string{"np"},
	Short:   "Show currently playing track",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		np, err := ops.GetNowPlaying(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, np, func() {
			status := "Paused"
			if np.IsPlaying {
				status = "Playing"
			}
			fmt.Printf("Status:   %s\n", status)
			fmt.Printf("Track:    %s\n", np.Track)
			fmt.Printf("Artist:   %s\n", np.Artist)
			fmt.Printf("Album:    %s\n", np.Album)
			fmt.Printf("Progress: %s / %s\n", np.Progress, np.Duration)
			fmt.Printf("Device:   %s\n", np.Device)
			fmt.Printf("Volume:   %d%%\n", np.Volume)
		})
		return nil
	},
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer playback to a device",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		play, _ := cmd.Flags().GetBool("play")
		if err := ops.TransferPlayback(c, deviceID, play); err != nil {
			return err
		}
		fmt.Println("Playback transferred")
		return nil
	},
}
