package cmd

import (
	"fmt"

	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queueCmd)
	queueCmd.AddCommand(queueListCmd, queueAddCmd)

	queueAddCmd.Flags().String("uri", "", "Spotify URI to add")
	queueAddCmd.MarkFlagRequired("uri")
	queueAddCmd.Flags().String("device", "", "Target device ID or name")
}

var queueCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"q"},
	Short:   "Queue commands",
	RunE:    queueListCmd.RunE,
}

var queueListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show playback queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		q, err := ops.GetQueue(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, q, func() {
			fmt.Printf("Now playing: %s - %s\n\n",
				q.NowPlaying.Name, q.NowPlaying.Artist)

			if len(q.Queue) == 0 {
				fmt.Println("Queue is empty")
				return
			}

			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "#\tTRACK\tARTIST\tDURATION")
			for i, t := range q.Queue {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, t.Name, t.Artist, t.Duration)
			}
			w.Flush()
		})
		return nil
	},
}

var queueAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a track to the queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		uri, _ := cmd.Flags().GetString("uri")
		deviceID, err := ops.ResolveDeviceID(c, resolveDevice(cmd))
		if err != nil {
			return err
		}
		if err := ops.AddToQueue(c, uri, deviceID); err != nil {
			return err
		}
		fmt.Println("Added to queue")
		return nil
	},
}
