---
page_title: "spotify_track Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Data Source `spotify_track`



## Example Usage

```terraform
data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"

  ## Computed
  # name = "Overkill"
  # artists = ["0qPGd8tOMHlFZt8EA1uLFY"]
  # album = "64ey3KHg3uepidKmJrb4ka"
}

data "spotify_track" "blackwater" {
  spotify_id = "4lE6N1E0L8CssgKEUCgdbA"

  ## Computed
  # name = "Blackwater"
  # artists = ["0qPGd8tOMHlFZt8EA1uLFY"]
  # album = "1AUS845POFhV3oDytPImEZ"
}
```

## Schema

### Optional

- **id** (String) The ID of this resource.
- **spotify_id** (String) Spotify ID of the track
- **url** (String) Spotify URL of the track

### Read-only

- **album** (String) The spotify ID of the album
- **artists** (List of String) The spotify IDs of the artists
- **name** (String) The Name of the track


