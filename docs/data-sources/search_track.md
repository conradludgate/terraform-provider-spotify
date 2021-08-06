---
page_title: "spotify_search_track Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Data Source `spotify_search_track`



## Example Usage

```terraform
resource "spotify_playlist" "ariana_grande" {
  name        = "My Ariana Grande Playlist"

  tracks = flatten([
    spotify_search_track.ariana_grande.tracks[*].id,
  ])
}

data "spotify_search_track" "ariana_grande" {
  artists = ["Ariana Grande"]
  limit = 10
}
```

## Schema

### Optional

- **album** (String) Name of the album
- **artist** (String) Name of the artist
- **explicit** (Boolean) Filter to allow explicit tracks
- **id** (String) The ID of this resource.
- **limit** (Number)
- **name** (String) Name of the track
- **year** (String) Year of release

### Read-only

- **tracks** (List of Object) List of tracks found (see [below for nested schema](#nestedatt--tracks))

<a id="nestedatt--tracks"></a>
### Nested Schema for `tracks`

Read-only:

- **album** (String)
- **artists** (List of String)
- **id** (String)
- **name** (String)


