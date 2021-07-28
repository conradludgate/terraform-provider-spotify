resource "spotify_playlist" "ariana_grande" {
  name        = "My Ariana Grande Playlist"

  tracks = flatten([
    spotify_search_track.ariana_grande.tracks[*].id,
    spotify_search_track.break_free.track.id
  ])
}

data "spotify_search_track" "ariana_grande" {
  artists = ["Ariana Grande"]
  limit = 10
}

data "spotify_search_track" "break_free" {
  artists = ["Ariana Grande", "Zedd"]
  name = "Break Free"
  limit = 1
}
