package cmd

import (
	"fmt"
	"strings"

	"github.com/jrogala/spotify-cli/internal/cmdutil"
	"github.com/jrogala/spotify-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().String("type", "track,album,artist", "Types to search (track,album,artist,playlist)")
	searchCmd.Flags().Int("limit", 10, "Max results per type")
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Spotify",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		query := strings.Join(args, " ")
		types, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := ops.Search(c, query, types, limit)
		if err != nil {
			return err
		}

		cmdutil.Render(cmd, result, func() {
			if len(result.Tracks) > 0 {
				fmt.Println("TRACKS")
				w := cmdutil.NewTabWriter()
				fmt.Fprintln(w, "  URI\tNAME\tARTIST\tALBUM\tDURATION")
				for _, t := range result.Tracks {
					fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n",
						t.URI, t.Name, t.Artist, t.Album, t.Duration)
				}
				w.Flush()
				fmt.Println()
			}

			if len(result.Albums) > 0 {
				fmt.Println("ALBUMS")
				w := cmdutil.NewTabWriter()
				fmt.Fprintln(w, "  URI\tNAME\tARTIST\tRELEASED\tTRACKS")
				for _, a := range result.Albums {
					fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%d\n",
						a.URI, a.Name, a.Artist, a.ReleaseDate, a.TotalTracks)
				}
				w.Flush()
				fmt.Println()
			}

			if len(result.Artists) > 0 {
				fmt.Println("ARTISTS")
				w := cmdutil.NewTabWriter()
				fmt.Fprintln(w, "  URI\tNAME")
				for _, a := range result.Artists {
					fmt.Fprintf(w, "  %s\t%s\n", a.URI, a.Name)
				}
				w.Flush()
				fmt.Println()
			}
		})
		return nil
	},
}
