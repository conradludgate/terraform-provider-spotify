---
page_title: "spotify Provider"
subcategory: ""
description: |-
  
---

# spotify Provider





## Schema

### Required

- **auth_code** (String) An authorization code required to exchange into an access token

### Optional

- **client_id** (String) The client ID used to exchange the auth code into an access token
- **client_secret** (String) The client secret used to convert the auth code into an access token. See https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow
- **code_verifier** (String) The code verifier value to exchange the auth code into an access token. See https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow-with-proof-key-for-code-exchange-pkce
- **redirect_uri** (String) The URI that spotify redirects to when authorizing and creating the auth_code
