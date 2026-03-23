// Package cmdutil provides shared helpers for CLI commands.
package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jrogala/spotify-cli/client"
	"github.com/jrogala/spotify-cli/config"
	"github.com/spf13/cobra"
)

// NewClient loads the token, refreshes if needed, and returns a Spotify client.
func NewClient() (*client.Client, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, err
	}

	if token.IsExpired() {
		resp, err := client.RefreshAccessToken(config.ClientID(), token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("refreshing token: %w", err)
		}
		token.AccessToken = resp.AccessToken
		token.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
		if resp.RefreshToken != "" {
			token.RefreshToken = resp.RefreshToken
		}
		if err := config.SaveToken(token); err != nil {
			return nil, fmt.Errorf("saving refreshed token: %w", err)
		}
	}

	return client.New(token.AccessToken), nil
}

// PrintJSON encodes v as indented JSON to stdout.
func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// IsJSON returns true if the --json persistent flag is set.
func IsJSON(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}

// Render outputs data as JSON if --json is set, otherwise calls tableFunc.
func Render(cmd *cobra.Command, data any, tableFunc func()) {
	if IsJSON(cmd) {
		PrintJSON(data)
		return
	}
	tableFunc()
}

// NewTabWriter creates a standard tabwriter for table output.
func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// ExitErr prints an error and exits.
func ExitErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
