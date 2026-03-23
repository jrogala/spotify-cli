package ops

import (
	"github.com/jrogala/spotify-cli/client"
)

// QueueEntry holds a queue item.
type QueueEntry struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
	URI      string `json:"uri"`
}

// QueueInfo holds the full queue state.
type QueueInfo struct {
	NowPlaying QueueEntry   `json:"now_playing"`
	Queue      []QueueEntry `json:"queue"`
}

func trackToEntry(t client.Track) QueueEntry {
	return QueueEntry{
		Name:     t.Name,
		Artist:   client.ArtistNames(t.Artists),
		Duration: client.FormatDuration(t.DurationMs),
		URI:      t.URI,
	}
}

// GetQueue returns the playback queue.
func GetQueue(c *client.Client) (*QueueInfo, error) {
	q, err := c.GetQueue()
	if err != nil {
		return nil, err
	}
	entries := make([]QueueEntry, len(q.Queue))
	for i, t := range q.Queue {
		entries[i] = trackToEntry(t)
	}
	return &QueueInfo{
		NowPlaying: trackToEntry(q.CurrentlyPlaying),
		Queue:      entries,
	}, nil
}

// AddToQueue adds a track to the queue.
func AddToQueue(c *client.Client, uri, deviceID string) error {
	return c.AddToQueue(uri, deviceID)
}
