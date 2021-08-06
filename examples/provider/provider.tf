provider "spotify" {
  api_key = var.spotify_api_key
}

# See https://github.com/conradludgate/terraform-provider-spotify#how-to-use
# for how to get an api key
variable "spotify_api_key" {
  type = string
}
