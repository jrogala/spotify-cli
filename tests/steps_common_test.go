package tests

import (
	"encoding/json"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/jrogala/spotify-cli/pkg/ops"
)

func registerSteps(sc *godog.ScenarioContext, ctx *scenarioCtx) {
	// ── Playback ──
	sc.Step(`^the Spotify API returns a playback state with track "([^"]*)" by "([^"]*)" on album "([^"]*)"$`, func(track, artist, album string) {
		ctx.mock.On("GET", "/me/player", 200, fmt.Sprintf(`{
			"is_playing": true,
			"progress_ms": 60000,
			"item": {
				"name": %q,
				"artists": [{"name": %q}],
				"album": {"name": %q},
				"duration_ms": 240000,
				"uri": "spotify:track:123"
			},
			"device": {"name": "TestDevice", "type": "Computer", "volume_percent": 75, "id": "dev1"},
			"shuffle_state": false,
			"repeat_state": "off"
		}`, track, artist, album))
	})

	sc.Step(`^I request the now-playing state$`, func() {
		np, err := ops.GetNowPlaying(ctx.client)
		ctx.err = err
		ctx.result = np
	})

	sc.Step(`^I should see track "([^"]*)"$`, func(expected string) error {
		np := ctx.result.(*ops.NowPlaying)
		if np.Track != expected {
			return fmt.Errorf("expected track %q, got %q", expected, np.Track)
		}
		return nil
	})

	sc.Step(`^I should see artist "([^"]*)"$`, func(expected string) error {
		np := ctx.result.(*ops.NowPlaying)
		if np.Artist != expected {
			return fmt.Errorf("expected artist %q, got %q", expected, np.Artist)
		}
		return nil
	})

	sc.Step(`^the Spotify API accepts play commands$`, func() {
		ctx.mock.On("PUT", "/me/player/play", 204, "")
	})

	sc.Step(`^I send a play command$`, func() {
		ctx.err = ops.Play(ctx.client, "", "", nil)
	})

	sc.Step(`^the command should succeed$`, func() error {
		if ctx.err != nil {
			return fmt.Errorf("expected success, got: %v", ctx.err)
		}
		return nil
	})

	sc.Step(`^the Spotify API accepts pause commands$`, func() {
		ctx.mock.On("PUT", "/me/player/pause", 204, "")
	})

	sc.Step(`^I send a pause command$`, func() {
		ctx.err = ops.Pause(ctx.client, "")
	})

	sc.Step(`^the Spotify API accepts next commands$`, func() {
		ctx.mock.On("POST", "/me/player/next", 204, "")
	})

	sc.Step(`^I send a next command$`, func() {
		ctx.err = ops.Next(ctx.client, "")
	})

	sc.Step(`^the Spotify API accepts previous commands$`, func() {
		ctx.mock.On("POST", "/me/player/previous", 204, "")
	})

	sc.Step(`^I send a previous command$`, func() {
		ctx.err = ops.Previous(ctx.client, "")
	})

	// ── Devices ──
	sc.Step(`^the Spotify API returns devices:$`, func(table *godog.Table) {
		var devices []map[string]any
		headers := table.Rows[0].Cells
		for _, row := range table.Rows[1:] {
			d := map[string]any{}
			for i, cell := range row.Cells {
				key := headers[i].Value
				switch key {
				case "is_active":
					d[key] = cell.Value == "true"
				case "volume_percent":
					var v int
					fmt.Sscanf(cell.Value, "%d", &v)
					d[key] = v
				default:
					d[key] = cell.Value
				}
			}
			devices = append(devices, d)
		}
		body, _ := json.Marshal(map[string]any{"devices": devices})
		ctx.mock.On("GET", "/me/player/devices", 200, string(body))
	})

	sc.Step(`^I request the device list$`, func() {
		result, err := ops.ListDevices(ctx.client)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) devices?$`, func(count int) error {
		devices := ctx.result.([]ops.DeviceInfo)
		if len(devices) != count {
			return fmt.Errorf("expected %d devices, got %d", count, len(devices))
		}
		return nil
	})

	sc.Step(`^device "([^"]*)" should be active$`, func(name string) error {
		devices := ctx.result.([]ops.DeviceInfo)
		for _, d := range devices {
			if d.Name == name && d.IsActive {
				return nil
			}
		}
		return fmt.Errorf("device %q is not active", name)
	})

	// ── Volume ──
	sc.Step(`^the Spotify API returns volume (\d+) on device "([^"]*)"$`, func(vol int, device string) {
		ctx.mock.On("GET", "/me/player", 200, fmt.Sprintf(`{
			"is_playing": true,
			"progress_ms": 0,
			"item": {"name": "", "artists": [], "album": {"name": ""}, "duration_ms": 0, "uri": ""},
			"device": {"name": %q, "volume_percent": %d, "id": "dev1"},
			"shuffle_state": false,
			"repeat_state": "off"
		}`, device, vol))
	})

	sc.Step(`^I request the current volume$`, func() {
		result, err := ops.GetVolume(ctx.client)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see volume (\d+)$`, func(expected int) error {
		vol := ctx.result.(*ops.VolumeInfo)
		if vol.Level != expected {
			return fmt.Errorf("expected volume %d, got %d", expected, vol.Level)
		}
		return nil
	})

	sc.Step(`^the Spotify API accepts volume commands$`, func() {
		ctx.mock.On("PUT", "/me/player/volume", 204, "")
	})

	sc.Step(`^I set the volume to (\d+)$`, func(level int) {
		ctx.err = ops.SetVolume(ctx.client, "", level)
	})

	// ── Queue ──
	sc.Step(`^the Spotify API returns a queue with current track "([^"]*)" and (\d+) queued tracks$`, func(current string, count int) {
		queue := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			queue[i] = map[string]any{
				"name":        fmt.Sprintf("Track %d", i+1),
				"artists":     []map[string]string{{"name": "Artist"}},
				"album":       map[string]string{"name": "Album"},
				"duration_ms": 200000,
				"uri":         fmt.Sprintf("spotify:track:q%d", i+1),
			}
		}
		body, _ := json.Marshal(map[string]any{
			"currently_playing": map[string]any{
				"name":        current,
				"artists":     []map[string]string{{"name": "Current Artist"}},
				"album":       map[string]string{"name": "Current Album"},
				"duration_ms": 180000,
				"uri":         "spotify:track:current",
			},
			"queue": queue,
		})
		ctx.mock.On("GET", "/me/player/queue", 200, string(body))
	})

	sc.Step(`^I request the queue$`, func() {
		result, err := ops.GetQueue(ctx.client)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^the now-playing track should be "([^"]*)"$`, func(expected string) error {
		q := ctx.result.(*ops.QueueInfo)
		if q.NowPlaying.Name != expected {
			return fmt.Errorf("expected now-playing %q, got %q", expected, q.NowPlaying.Name)
		}
		return nil
	})

	sc.Step(`^the queue should have (\d+) tracks$`, func(count int) error {
		q := ctx.result.(*ops.QueueInfo)
		if len(q.Queue) != count {
			return fmt.Errorf("expected %d queued tracks, got %d", count, len(q.Queue))
		}
		return nil
	})

	sc.Step(`^the Spotify API accepts add-to-queue commands$`, func() {
		ctx.mock.On("POST", "/me/player/queue", 204, "")
	})

	sc.Step(`^I add "([^"]*)" to the queue$`, func(uri string) {
		ctx.err = ops.AddToQueue(ctx.client, uri, "")
	})

	// ── Search ──
	sc.Step(`^the Spotify API returns search results for "([^"]*)" with (\d+) tracks and (\d+) albums and (\d+) artists$`, func(query string, nTracks, nAlbums, nArtists int) {
		tracks := make([]map[string]any, nTracks)
		for i := 0; i < nTracks; i++ {
			tracks[i] = map[string]any{
				"name":        fmt.Sprintf("Track Result %d", i+1),
				"artists":     []map[string]string{{"name": "Search Artist"}},
				"album":       map[string]string{"name": "Search Album"},
				"duration_ms": 210000,
				"uri":         fmt.Sprintf("spotify:track:s%d", i+1),
			}
		}
		albums := make([]map[string]any, nAlbums)
		for i := 0; i < nAlbums; i++ {
			albums[i] = map[string]any{
				"name":         fmt.Sprintf("Album Result %d", i+1),
				"artists":      []map[string]string{{"name": "Search Artist"}},
				"release_date": "2024-01-01",
				"total_tracks": 12,
				"uri":          fmt.Sprintf("spotify:album:s%d", i+1),
			}
		}
		artists := make([]map[string]any, nArtists)
		for i := 0; i < nArtists; i++ {
			artists[i] = map[string]any{
				"name": fmt.Sprintf("Artist Result %d", i+1),
				"uri":  fmt.Sprintf("spotify:artist:s%d", i+1),
			}
		}
		body, _ := json.Marshal(map[string]any{
			"tracks":  map[string]any{"items": tracks},
			"albums":  map[string]any{"items": albums},
			"artists": map[string]any{"items": artists},
		})
		ctx.mock.On("GET", "/search", 200, string(body))
	})

	sc.Step(`^I search for "([^"]*)" with type "([^"]*)"$`, func(query, types string) {
		result, err := ops.Search(ctx.client, query, types, 10)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) track results$`, func(count int) error {
		sr := ctx.result.(*ops.SearchResults)
		if len(sr.Tracks) != count {
			return fmt.Errorf("expected %d tracks, got %d", count, len(sr.Tracks))
		}
		return nil
	})

	sc.Step(`^I should see (\d+) album results$`, func(count int) error {
		sr := ctx.result.(*ops.SearchResults)
		if len(sr.Albums) != count {
			return fmt.Errorf("expected %d albums, got %d", count, len(sr.Albums))
		}
		return nil
	})

	sc.Step(`^I should see (\d+) artist results$`, func(count int) error {
		sr := ctx.result.(*ops.SearchResults)
		if len(sr.Artists) != count {
			return fmt.Errorf("expected %d artists, got %d", count, len(sr.Artists))
		}
		return nil
	})

	// ── Library ──
	sc.Step(`^the Spotify API returns (\d+) top tracks$`, func(count int) {
		items := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			items[i] = map[string]any{
				"name":    fmt.Sprintf("Top Track %d", i+1),
				"artists": []map[string]string{{"name": "Top Artist"}},
				"album":   map[string]string{"name": "Top Album"},
				"uri":     fmt.Sprintf("spotify:track:t%d", i+1),
			}
		}
		body, _ := json.Marshal(map[string]any{"items": items, "total": count})
		ctx.mock.On("GET", "/me/top/tracks", 200, string(body))
	})

	sc.Step(`^I request my top tracks$`, func() {
		result, err := ops.GetTopTracks(ctx.client, "medium_term", 20)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) top tracks$`, func(count int) error {
		tracks := ctx.result.([]ops.TopTrack)
		if len(tracks) != count {
			return fmt.Errorf("expected %d top tracks, got %d", count, len(tracks))
		}
		return nil
	})

	sc.Step(`^the Spotify API returns (\d+) top artists$`, func(count int) {
		items := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			items[i] = map[string]any{
				"name": fmt.Sprintf("Top Artist %d", i+1),
				"uri":  fmt.Sprintf("spotify:artist:a%d", i+1),
			}
		}
		body, _ := json.Marshal(map[string]any{"items": items, "total": count})
		ctx.mock.On("GET", "/me/top/artists", 200, string(body))
	})

	sc.Step(`^I request my top artists$`, func() {
		result, err := ops.GetTopArtists(ctx.client, "medium_term", 20)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) top artists$`, func(count int) error {
		artists := ctx.result.([]ops.TopArtist)
		if len(artists) != count {
			return fmt.Errorf("expected %d top artists, got %d", count, len(artists))
		}
		return nil
	})

	sc.Step(`^the Spotify API returns (\d+) liked tracks with total (\d+)$`, func(count, total int) {
		items := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			items[i] = map[string]any{
				"added_at": "2024-01-01T00:00:00Z",
				"track": map[string]any{
					"name":    fmt.Sprintf("Liked Track %d", i+1),
					"artists": []map[string]string{{"name": "Liked Artist"}},
					"album":   map[string]string{"name": "Liked Album"},
					"uri":     fmt.Sprintf("spotify:track:l%d", i+1),
				},
			}
		}
		body, _ := json.Marshal(map[string]any{"items": items, "total": total})
		ctx.mock.On("GET", "/me/tracks", 200, string(body))
	})

	sc.Step(`^I request my liked tracks$`, func() {
		result, err := ops.GetLikedTracks(ctx.client, 20, 0)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) liked tracks$`, func(count int) error {
		lr := ctx.result.(*ops.LikedResult)
		if len(lr.Items) != count {
			return fmt.Errorf("expected %d liked tracks, got %d", count, len(lr.Items))
		}
		return nil
	})

	sc.Step(`^the total should be (\d+)$`, func(total int) error {
		lr := ctx.result.(*ops.LikedResult)
		if lr.Total != total {
			return fmt.Errorf("expected total %d, got %d", total, lr.Total)
		}
		return nil
	})

	sc.Step(`^the Spotify API returns (\d+) recently played tracks$`, func(count int) {
		items := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			items[i] = map[string]any{
				"played_at": "2024-01-01T12:00:00Z",
				"track": map[string]any{
					"name":    fmt.Sprintf("Recent Track %d", i+1),
					"artists": []map[string]string{{"name": "Recent Artist"}},
					"album":   map[string]string{"name": "Recent Album"},
					"uri":     fmt.Sprintf("spotify:track:r%d", i+1),
				},
			}
		}
		body, _ := json.Marshal(map[string]any{"items": items})
		ctx.mock.On("GET", "/me/player/recently-played", 200, string(body))
	})

	sc.Step(`^I request my recently played tracks$`, func() {
		result, err := ops.GetRecentlyPlayed(ctx.client, 20)
		ctx.err = err
		ctx.result = result
	})

	sc.Step(`^I should see (\d+) recent tracks$`, func(count int) error {
		tracks := ctx.result.([]ops.RecentTrack)
		if len(tracks) != count {
			return fmt.Errorf("expected %d recent tracks, got %d", count, len(tracks))
		}
		return nil
	})
}
