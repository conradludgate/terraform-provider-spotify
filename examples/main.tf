terraform {
  required_providers {
    spotify = {
      version = "~> 0.1.2"
      source  = "conradludgate/spotify"
    }
  }
}

variable "spotify_auth_code" {
  type = string
}

variable "spotify_code_verifier" {
  type = string
}

spotify {
  auth_code     = var.spotify_auth_code
  code_verifier = var.spotify_code_verifier
}
