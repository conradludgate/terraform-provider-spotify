---
page_title: "spotify_playlist Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  Resource to manage a spotify playlist.
---

# Resource `spotify_playlist`

Resource to manage a spotify playlist.

## Example Usage

```terraform
resource "spotify_playlist" "playlist" {
  name        = "My playlist"
  description = "My playlist is so awesome"
  public      = false

  tracks = [
    data.spotify_track.overkill.id,
    data.spotify_track.blackwater.id,
    data.spotify_track.snowblind.id,
  ]
}

data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
}
data "spotify_track" "blackwater" {
  url = "https://open.spotify.com/track/4lE6N1E0L8CssgKEUCgdbA"
}
data "spotify_track" "snowblind" {
  url = "https://open.spotify.com/track/7FCG2wIYG1XvGRUMACC2cD"
}
```

## Schema

### Required

- **name** (String) The name of the resulting playlist
- **tracks** (Set of String) A set of tracks for the playlist to contain

### Optional

- **description** (String) The description of the resulting playlist
- **id** (String) The ID of this resource.
- **public** (Boolean) Whether the playlist can be accessed publically

### Read-only

- **snapshot_id** (String)

## Import

Import is supported using the following syntax:

```shell
# Using the playlist ID
# https://open.spotify.com/playlist/37i9dQZF1DWVs8I62NcHks (a playlist share link)
#                                   ^^^^^^^^^^^^^^^^^^^^^^
terraform import spotify_playlist.example 37i9dQZF1DWVs8I62NcHks
```
