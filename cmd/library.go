package cmd

import (
	"fmt"

	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(topCmd, likedCmd, recentCmd)

	topCmd.Flags().String("type", "tracks", "Type: tracks or artists")
	topCmd.Flags().String("range", "medium_term", "Time range: short_term (4w), medium_term (6m), long_term (years)")
	topCmd.Flags().Int("limit", 20, "Max results")

	likedCmd.Flags().Int("limit", 20, "Max results")
	likedCmd.Flags().Int("offset", 0, "Offset")

	recentCmd.Flags().Int("limit", 20, "Max results")
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show your top tracks or artists",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		typ, _ := cmd.Flags().GetString("type")
		timeRange, _ := cmd.Flags().GetString("range")
		limit, _ := cmd.Flags().GetInt("limit")

		if typ == "artists" {
			artists, err := ops.GetTopArtists(c, timeRange, limit)
			if err != nil {
				return err
			}
			cmdutil.Render(cmd, artists, func() {
				w := cmdutil.NewTabWriter()
				fmt.Fprintln(w, "#\tARTIST\tURI")
				for i, a := range artists {
					fmt.Fprintf(w, "%d\t%s\t%s\n", i+1, a.Name, a.URI)
				}
				w.Flush()
			})
		} else {
			tracks, err := ops.GetTopTracks(c, timeRange, limit)
			if err != nil {
				return err
			}
			cmdutil.Render(cmd, tracks, func() {
				w := cmdutil.NewTabWriter()
				fmt.Fprintln(w, "#\tTRACK\tARTIST\tALBUM\tURI")
				for i, t := range tracks {
					fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
						i+1, t.Name, t.Artist, t.Album, t.URI)
				}
				w.Flush()
			})
		}
		return nil
	},
}

var likedCmd = &cobra.Command{
	Use:     "liked",
	Aliases: []string{"saved"},
	Short:   "Show your liked/saved tracks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		result, err := ops.GetLikedTracks(c, limit, offset)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, result, func() {
			fmt.Printf("Total liked tracks: %d\n\n", result.Total)
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "#\tTRACK\tARTIST\tALBUM")
			for i, t := range result.Items {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					offset+i+1, t.Name, t.Artist, t.Album)
			}
			w.Flush()
		})
		return nil
	},
}

var recentCmd = &cobra.Command{
	Use:     "recent",
	Aliases: []string{"history"},
	Short:   "Show recently played tracks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")

		history, err := ops.GetRecentlyPlayed(c, limit)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, history, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "PLAYED AT\tTRACK\tARTIST")
			for _, h := range history {
				played := h.PlayedAt
				if len(played) > 16 {
					played = played[:16]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", played, h.Name, h.Artist)
			}
			w.Flush()
		})
		return nil
	},
}
