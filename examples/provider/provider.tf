provider "spotify" {
  api_key = var.spotify_api_key
}

# See https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy
# for how to get an api key
variable "spotify_api_key" {
  type = string
}
