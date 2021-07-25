---
page_title: "spotify_playlist Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  Resource to manage a spotify playlist.
---

# Resource `spotify_playlist`

Resource to manage a spotify playlist.



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

A Spotify playlist can be imported using the playlist id, e.g.,

```
$ terraform import spotify_playlist.example 37i9dQZF1DWVs8I62NcHks
```

The playlist id is part of the URL used by the Spotify Web Player, e.g.,

```
https://open.spotify.com/playlist/37i9dQZF1DWVs8I62NcHks
```
