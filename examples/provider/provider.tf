provider "spotify" {
  api_key = var.spotify_api_key
}

variable "spotify_api_key" {
  type = string
}
