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

# See https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy
# for how to get an api key
variable "spotify_api_key" {
  type = string
}
```

## Schema

### Required

- **api_key** (String) Oauth2 Proxy API Key

### Optional

- **auth_server** (String) Oauth2 Proxy URL
- **token_id** (String) Oauth2 Proxy token ID
- **username** (String) Oauth2 Proxy username
