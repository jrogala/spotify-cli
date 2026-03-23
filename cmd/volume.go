package cmd

import (
	"fmt"
	"strconv"

	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.Flags().String("device", "", "Target device ID or name")
}

var volumeCmd = &cobra.Command{
	Use:   "volume [level]",
	Short: "Get or set volume (0-100)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			vol, err := ops.GetVolume(c)
			if err != nil {
				return err
			}
			cmdutil.Render(cmd, vol, func() {
				fmt.Printf("%d%% (%s)\n", vol.Level, vol.Device)
			})
			return nil
		}

		level, err := strconv.Atoi(args[0])
		if err != nil || level < 0 || level > 100 {
			return fmt.Errorf("volume must be 0-100")
		}
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		if err := ops.SetVolume(c, deviceID, level); err != nil {
			return err
		}
		fmt.Printf("Volume set to %d%%\n", level)
		return nil
	},
}
