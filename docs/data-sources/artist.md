---
page_title: "spotify_artist Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Data Source `spotify_artist`



## Example Usage

```terraform
data "spotify_artist" "overdrive" {
  url = "https://open.spotify.com/artist/3PALZKWkpwjRvBsRmhlVSS"

  ## Computed
  # name = "Gunship"
}

data "spotify_artist" "wolfclub" {
  spotify_id = "4dCDYKtFTMnKCI9PvEwMQX"

  ## Computed
  # name = "W O L F C L U B"
}
```

## Schema

### Optional

- **spotify_id** (String) Spotify ID of the artist
- **url** (String) Spotify URL of the artist

### Read-only

- **name** (String) The Name of the artist


