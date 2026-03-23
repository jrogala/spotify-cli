package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jrogala/spotify-cli/client"
	"github.com/jrogala/spotify-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd, authLogoutCmd, authStatusCmd)
	authLoginCmd.Flags().Bool("manual", false, "Manually paste the callback URL (for sandboxed browsers)")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Spotify via browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID := config.ClientID()
		redirectURI := config.RedirectURI()

		verifier, err := client.GenerateCodeVerifier()
		if err != nil {
			return err
		}
		challenge := client.CodeChallenge(verifier)
		state, err := client.GenerateState()
		if err != nil {
			return err
		}

		authURL := client.AuthorizeURL(clientID, redirectURI, challenge, state)

		manual, _ := cmd.Flags().GetBool("manual")

		fmt.Println("Open this URL in your browser:")
		fmt.Println()
		fmt.Println(authURL)
		fmt.Println()

		var code string

		if manual {
			fmt.Println("After authorizing, paste the full callback URL here:")
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return fmt.Errorf("no input")
			}
			callbackURL := strings.TrimSpace(scanner.Text())
			parsed, err := url.Parse(callbackURL)
			if err != nil {
				return fmt.Errorf("invalid URL: %w", err)
			}
			if parsed.Query().Get("state") != state {
				return fmt.Errorf("state mismatch")
			}
			if e := parsed.Query().Get("error"); e != "" {
				return fmt.Errorf("authorization denied: %s", e)
			}
			code = parsed.Query().Get("code")
		} else {
			codeCh := make(chan string, 1)
			errCh := make(chan error, 1)

			mux := http.NewServeMux()
			mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("state") != state {
					http.Error(w, "Invalid state", http.StatusBadRequest)
					errCh <- fmt.Errorf("state mismatch")
					return
				}
				if e := r.URL.Query().Get("error"); e != "" {
					http.Error(w, "Authorization denied", http.StatusForbidden)
					errCh <- fmt.Errorf("authorization denied: %s", e)
					return
				}
				c := r.URL.Query().Get("code")
				fmt.Fprintf(w, "<html><body><h2>Authenticated! You can close this tab.</h2></body></html>")
				codeCh <- c
			})

			listener, err := net.Listen("tcp", "127.0.0.1:8888")
			if err != nil {
				return fmt.Errorf("failed to start callback server: %w", err)
			}

			server := &http.Server{Handler: mux}
			go server.Serve(listener)

			fmt.Println("Waiting for authorization...")

			select {
			case code = <-codeCh:
			case err := <-errCh:
				server.Shutdown(context.Background())
				return err
			case <-time.After(5 * time.Minute):
				server.Shutdown(context.Background())
				return fmt.Errorf("timed out waiting for authorization")
			}

			server.Shutdown(context.Background())
		}

		tokenResp, err := client.ExchangeCode(clientID, code, redirectURI, verifier)
		if err != nil {
			return err
		}

		token := &config.Token{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			ExpiresAt:    tokenResp.ExpiresAt(),
			Scope:        tokenResp.Scope,
		}
		if err := config.SaveToken(token); err != nil {
			return err
		}

		fmt.Println("Authenticated successfully!")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.DeleteToken(); err != nil {
			return err
		}
		fmt.Println("Logged out")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := config.LoadToken()
		if err != nil {
			fmt.Println("Not authenticated")
			return nil
		}
		if token.IsExpired() {
			fmt.Println("Token expired (will auto-refresh on next command)")
		} else {
			fmt.Printf("Authenticated (expires %s)\n", token.ExpiresAt.Format(time.RFC3339))
		}
		fmt.Printf("Scopes: %s\n", token.Scope)
		return nil
	},
}
