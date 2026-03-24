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
$ spotify devices
ID                                        NAME              TYPE     ACTIVE  VOLUME
abc123def456789012345678901234567890abcdef  Living Room       speaker  *       85%
def456789012345678901234567890abcdef123456  Kitchen Speaker   speaker          60%

$ spotify now-playing
Status:   Playing
Track:    Bohemian Rhapsody
Artist:   Queen
Album:    A Night at the Opera
Progress: 3:22 / 5:55
Device:   Living Room
Volume:   85%

$ spotify search "Bohemian Rhapsody"
TRACKS
  URI                           NAME               ARTIST  ALBUM                 DURATION
  spotify:track:xxx             Bohemian Rhapsody  Queen   A Night at the Opera  5:55
  spotify:track:yyy             Bohemian Rhapsody  Queen   Greatest Hits         5:55

$ spotify top --type tracks --range short
#   TRACK              ARTIST        ALBUM              URI
1   Blinding Lights    The Weeknd    After Hours        spotify:track:xxx
2   Shape of You       Ed Sheeran    ÷ (Divide)         spotify:track:yyy
```

## JSON Output

All commands support `--json` for machine-readable output.
