resource "spotify_playlist" "ariana_grande" {
  name        = "My Ariana Grande Playlist"

  tracks = flatten([
    spotify_search_track.ariana_grande.tracks[*].id,
  ])
}

data "spotify_search_track" "ariana_grande" {
  artist = "Ariana Grande"
  limit = 10
}
