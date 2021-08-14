---
page_title: "spotify_library_albums Resource - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Resource `spotify_library_albums`



## Example Usage

```terraform
resource "spotify_library_albums" "my_album" {
  albums = [
    data.spotify_album.only_in_dreams.id,
    data.spotify_album.the_promised_land.id,
  ]
}

data "spotify_album" "only_in_dreams" {
  spotify_id = "35axN2yrxRiycF2pA8mZaB"
}

data "spotify_album" "the_promised_land" {
  url = "https://open.spotify.com/album/3nRnJkUJYFfxcOGgU6LNci"
}
```

## Schema

### Required

- **albums** (Set of String) The list of track IDs to save to your 'liked albums'. *Note, if used incorrectly you may unlike all of your albums - use with caution*

### Optional

- **id** (String) The ID of this resource.


