package cmd

import (
	"fmt"

	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List available Spotify Connect devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		devices, err := ops.ListDevices(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, devices, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME\tTYPE\tACTIVE\tVOLUME")
			for _, d := range devices {
				active := ""
				if d.IsActive {
					active = "*"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d%%\n", d.ID, d.Name, d.Type, active, d.Volume)
			}
			w.Flush()
		})
		return nil
	},
}
