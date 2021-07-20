---
page_title: "spotify_search_track Data Source - terraform-provider-spotify"
subcategory: ""
description: |-
  
---

# Data Source `spotify_search_track`





## Schema

### Optional

- **album** (String)
- **artists** (List of String)
- **id** (String) The ID of this resource.
- **limit** (Number)
- **name** (String)
- **year** (String)
- **explicit** (Bool)

### Read-only

- **track** (Map of String)
- **tracks** (List of Object) (see [below for nested schema](#nestedatt--tracks))

<a id="nestedatt--tracks"></a>
### Nested Schema for `tracks`

Read-only:

- **album** (String)
- **artists** (List of String)
- **id** (String)
- **name** (String)


