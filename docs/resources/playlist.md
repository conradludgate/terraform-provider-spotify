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
- **tracks** (List of String) A list of tracks for the playlist to contain

### Optional

- **description** (String) The description of the resulting playlist
- **id** (String) The ID of this resource.
- **public** (Boolean) Whether the playlist can be accessed publically


