package client

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	authURL  = "https://accounts.spotify.com/authorize"
	tokenURL = "https://accounts.spotify.com/api/token"

	Scopes = "ugc-image-upload user-read-playback-state user-modify-playback-state user-read-currently-playing user-top-read user-library-read user-library-modify user-read-playback-position user-read-recently-played user-follow-read user-follow-modify playlist-read-private playlist-read-collaborative playlist-modify-public playlist-modify-private app-remote-control streaming"
)

// TokenResponse is the response from the Spotify token endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateCodeVerifier creates a random PKCE code verifier.
func GenerateCodeVerifier() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// CodeChallenge computes the S256 PKCE code challenge from a verifier.
func CodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// AuthorizeURL builds the Spotify authorization URL for PKCE flow.
func AuthorizeURL(clientID, redirectURI, codeChallenge, state string) string {
	v := url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"scope":                 {Scopes},
		"redirect_uri":          {redirectURI},
		"code_challenge_method": {"S256"},
		"code_challenge":        {codeChallenge},
		"state":                 {state},
	}
	return authURL + "?" + v.Encode()
}

// ExchangeCode exchanges an authorization code for tokens using PKCE.
func ExchangeCode(clientID, code, redirectURI, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"code_verifier": {codeVerifier},
	}
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("token exchange failed: %s (%s)", errResp.Error, errResp.Description)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	return &tokenResp, nil
}

// RefreshAccessToken refreshes an access token using a refresh token.
func RefreshAccessToken(clientID, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
	}
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token refresh: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("token refresh failed: %s (%s)", errResp.Error, errResp.Description)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	return &tokenResp, nil
}

// GenerateState creates a random state string for CSRF protection.
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ExpiresAt computes the token expiration time from ExpiresIn seconds.
func (t *TokenResponse) ExpiresAt() time.Time {
	return time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
}
