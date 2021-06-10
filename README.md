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

Add

The following to your terraform configuration

```tf
terraform {
  required_providers {
    spotify = {
      version = "~> 0.1.4"
      source  = "conradludgate/spotify"
    }
  }
}
```

## How to use

First, you need an instance of a spotify auth server running. This acts as a middleware between terraform and spotify to allow easy access to access tokens.
See [spotify_auth_proxy](/spotify_auth_proxy) to get started.

Once you have the server running, make note of the API Key it gives you.

Configure the terraform provider like so

```tf
variable "spotify_api_key" {
  type = string
}

provider "spotify" {
  api_key = var.spotify_api_key
}
```

## Example

```tf
resource "spotify_playlist" "playlist" {
  name        = "My playlist"
  description = "My playlist is so awesome"
  public      = false

  tracks = [
    data.spotify_track.overkill.id,
    data.spotify_track.blackwater.id,
    data.spotify_track.overkill.id,
    data.spotify_search_track.search.tracks[0].id,
  ]
}

data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
}
data "spotify_track" "blackwater" {
  spotify_id = "4lE6N1E0L8CssgKEUCgdbA"
}

data "spotify_search_track" "search" {
  name    = "Somebody Told Me"
  artists = ["The Killers"]
  album   = "Hot Fuss"
}

output "test" {
  value = data.spotify_search_track.search.tracks
}
```
