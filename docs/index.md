---
page_title: "spotify Provider"
subcategory: ""
description: |-
  
---

# spotify Provider



## Example Usage

```terraform
provider "spotify" {
  api_key = var.spotify_api_key
}

variable "spotify_api_key" {
  type = string
}
```

## Schema

### Required

- **api_key** (String) Spotify auth proxy API Key

### Optional

- **auth_server** (String) Spotify auth proxy URL
