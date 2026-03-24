# spotify-cli

CLI for Spotify playback control and library browsing.

## Install

Download a binary from the [latest release](https://github.com/jrogala/spotify-cli/releases/latest), or install with Go:

```bash
go install github.com/jrogala/spotify-cli@latest
```

## Setup

Authenticate via OAuth2:

```bash
spotify auth login
```

Config stored at `~/.config/spotify-cli/config.yaml`.

## Commands

| Command | Description |
|---|---|
| `play` | Start or resume playback |
| `pause` | Pause playback |
| `next` | Skip to next track |
| `prev` | Skip to previous track |
| `now-playing` | Show currently playing track |
| `volume` | Get or set volume (0-100) |
| `devices` | List available playback devices |
| `transfer` | Transfer playback to another device |
| `search` | Search for tracks, albums, or artists |
| `queue list` | Show the current play queue |
| `queue add` | Add a track to the queue |
| `top` | Show top tracks or artists |
| `liked` | List liked/saved tracks |
| `recent` | Show recently played tracks |
| `auth login` | Authenticate via OAuth2 |
| `auth logout` | Remove stored credentials |
| `auth status` | Show current auth status |

## Examples

```bash
# Start playback on a specific device
spotify play --device "Kitchen Speaker"

# Search and queue a track
spotify search "Bohemian Rhapsody"
spotify queue add --uri spotify:track:xxx

# Transfer playback to another device
spotify transfer --device "Living Room"

# Show your top tracks from the last month
spotify top --type tracks --range short
```

## JSON Output

All commands support `--json` for machine-readable output.
