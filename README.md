# spotify-cli

CLI for Spotify playback control and library browsing.

## Install

Download a binary from the [latest release](https://github.com/jrogala/spotify-cli/releases/latest), or install with Go:

```bash
go install github.com/jrogala/spotify-cli@latest
```

## Setup

1. Create an app at [developer.spotify.com/dashboard](https://developer.spotify.com/dashboard), set redirect URI to `http://127.0.0.1:8888/callback`
2. Add your client ID to `~/.config/spotify-cli/config.yaml`:
   ```yaml
   client_id: your_client_id
   ```
3. Run `spotify auth login`

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
ID                                        NAME                  TYPE      ACTIVE  VOLUME
434ab3de26b2f6be96d6666d2fe98e777163df42  bureau                Speaker   *       15%
4469fbb3e9f7ea3ba1ede58f4ba49da9c90994f8  Web Player (Firefox)  Computer          10%

$ spotify now-playing
Status:   Playing
Track:    Dean Town
Artist:   Vulfpeck
Album:    The Beautiful Game
Progress: 3:16 / 3:33
Device:   bureau
Volume:   15%

$ spotify search "Vulfpeck" --type artist
ARTISTS
  URI                                     NAME
  spotify:artist:7pXu47GoqSYRajmBCjxdD6  Vulfpeck
  spotify:artist:6xt9sJmmyYwWkJv8A6ssiU  Cory Wong
  spotify:artist:1JyLSGXC3aWzjY6ZdxvIXh  The Fearless Flyers

$ spotify top --type tracks --range short_term
#   TRACK                 ARTIST              ALBUM                                URI
1   トドメの一撃              Vaundy, Cory Wong   replica                              spotify:track:7sd09c...
2   Stay With Me          Cory Wong           Lost In The Wonder                   spotify:track:5dcr3Z...
3   Master Of Puppets     Metallica           Master Of Puppets (Remastered)       spotify:track:2MuWTI...
```

## JSON Output

All commands support `--json` for machine-readable output.
