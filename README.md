# terraform-provider-spotify

This is a terraform provider for managing your spotify playlists.
- [terraform-provider-spotify](#terraform-provider-spotify)
  - [Installation](#installation)
  - [How to use](#how-to-use)
  - [Example](#example)
  - [Documentation](#documentation)
    - [Resources](#resources)
      - [spotify_library](#spotify_library)
      - [spotify_playlist](#spotify_playlist)
    - [Data sources](#data-sources)
      - [spotify_search_track](#spotify_search_track)
  - [Todo](#todo)
    - [Playlist diff](#playlist-diff)
    - [More Datasources](#more-datasources)

## Installation

To install the provider, run
```sh
make install
```
This will build the binary and install it into the terraform plugins dir.

## How to use

To run `terraform plan` or `terraform apply`, you must have a valid access token.
See [spotify_auth](/tree/main/spotify_auth) for information about how to get an access token.

The provider will look for access tokens in the environment variable `SPOTIFY_ACCESS_TOKEN`.
The access token must be valid for the following scopes:
*   user-read-email
*   user-read-private
*   playlist-read-private
*   playlist-modify-private
*   playlist-modify-public
*   user-library-read
*   user-library-modify

Some of the scopes may be omitted based on the resources you use.

For more information, see
https://developer.spotify.com/documentation/general/guides/authorization-guide/#implicit-grant-flow
https://developer.spotify.com/documentation/general/guides/scopes/

## Example

```tf
# Creates a private playlist named "My playlist"
# and adds 2 tracks to it
resource "spotify_playlist" "playlist" {
    name = "My playlist"
    description = "My playlist is so awesome"
    public = false

    tracks = [
        data.spotify_search_track.overkill.id,
        data.spotify_search_track.blackwater.id,
    ]
}

# Searches spotify for "Overkill RIOT", returns the first track ID
data "spotify_search_track" "overkill" {
    name = "Overkill"

    artists = [
        "RIOT",
    ]
}

# Searches spotify for "Blackwater RIOT", returns the first track ID
data "spotify_search_track" "blackwater" {
    name = "Blackwater"

    artists = [
        "RIOT",
    ]
}
```

## Documentation

### Resources

#### spotify_library

Manage the tracks saved in the user's 'Liked Songs'

**Required scopes:**
*   [user-library-read](https://developer.spotify.com/documentation/general/guides/scopes/#user-library-read) - To read the tracks that are currently saved
*   [user-library-modify](https://developer.spotify.com/documentation/general/guides/scopes/#user-library-modify) - To save/remove tracks from the library

**Variables:**
*   **tracks**: *set[string]* - Set of the tracks IDs in the playlist

#### spotify_playlist

Manage the details and tracks in a playlist

**Required scopes:**
*   [user-read-private](https://developer.spotify.com/documentation/general/guides/scopes/#user-read-private) - To get the User ID associated with the access token
*   [playlist-read-private](https://developer.spotify.com/documentation/general/guides/scopes/#playlist-read-private) - To read the tracks in a private playlist (only needed if `public = false`)
*   [playlist-modify-private](https://developer.spotify.com/documentation/general/guides/scopes/#playlist-modify-private) - To create/update the tracks in a private playlist (only needed if `public = false`)
*   [playlist-modify-public](https://developer.spotify.com/documentation/general/guides/scopes/#playlist-modify-public) - To create/update the tracks in a public playlist (only needed if `public = true`)

**Variables:**
*   **name**: *string* - Name of the playlist
*   **description**: *string?* - Description of the playlist
*   **public**: *bool?* - Whether the playlist is public (default `true`)
*   **tracks**: *list[string]* - List of the tracks IDs in the playlist

**Computed:**
*   **id**: *string* - Playlist ID

### Data sources

#### spotify_search_track

Search for a track

**Required scopes:**
None

**Paramaters:**
*   **name**: *string* - Name of the track
*   **artists**: *list[string]* - List of the artists

**Results:**
*   **id**: *string* - ID of the first track found


## Todo

### Playlist diff
Currently, the playlists tracks don't perform a diff,
they instead just replace all the tracks in the playlist with the new set.
This is fine for playlists with less than a few hundred songs,
But for larger playlists, adding a single song will cause several large API requests.

For efficiency, a diff of the old and new track list should be calculated.
The exact number of API calls can be known before starting,
choosing the lesser will most likely be the better option.

### More Datasources
For my personal uses, I would like to have my 'Liked Songs' as a data source,
and not something I manage personally.

I have 3000 songs saved,
I could write a script to convert them into a terraform file,
but keeping it synced will be a nightmare.

So instead, I want it as a data source I can use in other terraform resources.

For example, a way to get a random playlist made out of my liked songs

```tf
provider "random_shuffle" {
    version = "2.3.0"
}

data "spotify_saved_songs" "my_songs" {

}

resource "random_shuffle" "shuffle_daily" {
  input = data.spotify_saved_songs.my_songs.tracks
  result_count = 100
  seed = formatdate("YYYY-MM-DD", "2018-01-02T23:12:01-08:00")
}

resource "spotify_playlist" "shuffled_playlist" {
    name = "Daily shuffle mix"
    description = "Check back daily for a new awesome mix of songs"
    public = true

    tracks = random_shuffle.shuffle_daily
}
```
