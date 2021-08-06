# terraform-provider-spotify

[![docs](https://img.shields.io/static/v1?label=docs&message=terraform&color=informational&style=for-the-badge)](https://registry.terraform.io/providers/conradludgate/spotify/latest/docs)

This is a terraform provider for managing your spotify playlists.

Featured tutorial - https://learn.hashicorp.com/tutorials/terraform/spotify-playlist

Featured interview - https://www.hashicorp.com/blog/build-your-summer-spotify-playlist-with-terraform

## Example

```tf
resource "spotify_playlist" "playlist" {
  name        = "My playlist"
  description = "My playlist is so awesome"
  public      = false

  tracks = flatten([
    data.spotify_track.overkill.id,
    data.spotify_track.blackwater.id,
    data.spotify_track.overkill.id,
    data.spotify_search_track.search.tracks[*].id,
  ])
}

data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
}
data "spotify_track" "blackwater" {
  spotify_id = "4lE6N1E0L8CssgKEUCgdbA"
}

data "spotify_search_track" "search" {
  name   = "Somebody Told Me"
  artist = "The Killers"
  album  = "Hot Fuss"
}

output "test" {
  value = data.spotify_search_track.search.tracks
}
```


## Installation

Add the following to your terraform configuration

```tf
terraform {
  required_providers {
    spotify = {
      source  = "conradludgate/spotify"
      version = "~> 0.2.0"
    }
  }
}
```

## How to use

First, you need an instance of a spotify oauth2 server running. This acts as a middleware between terraform and spotify to allow easy access to access tokens.

### Public proxy

For a simple way to manage your spotify oauth2 tokens is to use https://oauth2.conrad.cafe. ([source code](https://github.com/conradludgate/oauth2-proxy))

Register a new account, create a spotify token with the following scopes

* user-read-email
* user-read-private
* playlist-read-private
* playlist-modify-private
* playlist-modify-public
* user-library-read
* user-library-modify

Then take note of the token id in the URL and the API key that is shown on the page

Configure the terraform provider like so

```tf
provider "spotify" {
  auth_server = "https://oauth2.conrad.cafe"
  api_key = var.spotify_api_key
  username = "your username"
  token_id = "your token id"
}

variable "spotify_api_key" {
  type = string
}
```

### Self hosted

If you want a bit more control over your tokens, you can self host a simple instance of the oauth2 proxy designed specifically for this terraform provider

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
