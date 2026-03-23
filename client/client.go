package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://api.spotify.com/v1"

// Client communicates with the Spotify Web API.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// New creates a new Spotify API client with the given access token.
func New(accessToken string) *Client {
	return &Client{
		token:      accessToken,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
}

// NewWithBaseURL creates a client that talks to the given base URL (for testing).
func NewWithBaseURL(accessToken, baseURL string) *Client {
	return &Client{
		token:      accessToken,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// do sends an HTTP request to the Spotify API.
func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encoding request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseError(resp.Body)
	}

	return resp, nil
}

// doJSON sends a request and decodes the JSON response into result.
func (c *Client) doJSON(method, path string, body any, result any) error {
	resp, err := c.do(method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// --- Player Commands ---

func (c *Client) GetPlaybackState() (*PlaybackState, error) {
	var state PlaybackState
	if err := c.doJSON("GET", "/me/player", nil, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (c *Client) GetDevices() ([]Device, error) {
	var result struct {
		Devices []Device `json:"devices"`
	}
	if err := c.doJSON("GET", "/me/player/devices", nil, &result); err != nil {
		return nil, err
	}
	return result.Devices, nil
}

func (c *Client) Play(deviceID string, contextURI string, uris []string, offsetPosition int) error {
	body := map[string]any{}
	if contextURI != "" {
		body["context_uri"] = contextURI
	}
	if len(uris) > 0 {
		body["uris"] = uris
	}
	if offsetPosition >= 0 {
		body["offset"] = map[string]int{"position": offsetPosition}
	}

	path := "/me/player/play"
	if deviceID != "" {
		path += "?device_id=" + url.QueryEscape(deviceID)
	}

	var reqBody any
	if len(body) > 0 {
		reqBody = body
	}

	resp, err := c.do("PUT", path, reqBody)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) Pause(deviceID string) error {
	path := "/me/player/pause"
	if deviceID != "" {
		path += "?device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("PUT", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) Next(deviceID string) error {
	path := "/me/player/next"
	if deviceID != "" {
		path += "?device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("POST", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) Previous(deviceID string) error {
	path := "/me/player/previous"
	if deviceID != "" {
		path += "?device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("POST", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) SetVolume(deviceID string, percent int) error {
	path := "/me/player/volume?volume_percent=" + strconv.Itoa(percent)
	if deviceID != "" {
		path += "&device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("PUT", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) TransferPlayback(deviceID string, play bool) error {
	body := map[string]any{
		"device_ids": []string{deviceID},
		"play":       play,
	}
	resp, err := c.do("PUT", "/me/player", body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) GetQueue() (*QueueResponse, error) {
	var q QueueResponse
	if err := c.doJSON("GET", "/me/player/queue", nil, &q); err != nil {
		return nil, err
	}
	return &q, nil
}

func (c *Client) AddToQueue(uri, deviceID string) error {
	path := "/me/player/queue?uri=" + url.QueryEscape(uri)
	if deviceID != "" {
		path += "&device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("POST", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) Shuffle(deviceID string, state bool) error {
	path := "/me/player/shuffle?state=" + strconv.FormatBool(state)
	if deviceID != "" {
		path += "&device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("PUT", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) Repeat(deviceID string, state string) error {
	path := "/me/player/repeat?state=" + url.QueryEscape(state)
	if deviceID != "" {
		path += "&device_id=" + url.QueryEscape(deviceID)
	}
	resp, err := c.do("PUT", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// --- Search ---

func (c *Client) Search(query, types string, limit int) (*SearchResult, error) {
	path := "/search?q=" + url.QueryEscape(query) +
		"&type=" + url.QueryEscape(types) +
		"&limit=" + strconv.Itoa(limit)
	var result SearchResult
	if err := c.doJSON("GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- User Library ---

func (c *Client) GetTopTracks(timeRange string, limit int) ([]Track, error) {
	path := "/me/top/tracks?limit=" + strconv.Itoa(limit) + "&time_range=" + url.QueryEscape(timeRange)
	var paging Paging[Track]
	if err := c.doJSON("GET", path, nil, &paging); err != nil {
		return nil, err
	}
	return paging.Items, nil
}

func (c *Client) GetTopArtists(timeRange string, limit int) ([]Artist, error) {
	path := "/me/top/artists?limit=" + strconv.Itoa(limit) + "&time_range=" + url.QueryEscape(timeRange)
	var paging Paging[Artist]
	if err := c.doJSON("GET", path, nil, &paging); err != nil {
		return nil, err
	}
	return paging.Items, nil
}

func (c *Client) GetLikedTracks(limit, offset int) (*Paging[SavedTrack], error) {
	path := "/me/tracks?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	var paging Paging[SavedTrack]
	if err := c.doJSON("GET", path, nil, &paging); err != nil {
		return nil, err
	}
	return &paging, nil
}

func (c *Client) GetRecentlyPlayed(limit int) ([]PlayHistory, error) {
	path := "/me/player/recently-played?limit=" + strconv.Itoa(limit)
	var result struct {
		Items []PlayHistory `json:"items"`
	}
	if err := c.doJSON("GET", path, nil, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// --- Albums ---

func (c *Client) GetAlbumTracks(albumID string) ([]Track, error) {
	var paging Paging[Track]
	if err := c.doJSON("GET", "/albums/"+url.PathEscape(albumID)+"/tracks?limit=50", nil, &paging); err != nil {
		return nil, err
	}
	return paging.Items, nil
}

// artistNames returns a comma-separated string of artist names.
func ArtistNames(artists []Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

// FormatDuration formats milliseconds as m:ss.
func FormatDuration(ms int) string {
	s := ms / 1000
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
