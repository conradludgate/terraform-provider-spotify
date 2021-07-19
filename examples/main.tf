terraform {
  required_providers {
    spotify = {
      version = "~> 0.1.7"
      source  = "conradludgate/spotify"
    }
  }
}

variable "spotify_api_key" {
  type = string
}

provider "spotify" {
  api_key = var.spotify_api_key
}
