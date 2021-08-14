---
page_title: "spotify_library Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Resource `spotify_library`



## Example Usage

```terraform
resource "spotify_library" "my_library" {
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

- **tracks** (Set of String) The list of track IDs to save to your 'liked tracks'. *Note, if used incorrectly you may unlike all of your tracks - use with caution*

### Optional

- **id** (String) The ID of this resource.


