resource "spotify_library" "my_library" {
  tracks = [
    data.spotify_track.overkill.id,
    data.spotify_track.blackwater.id,
    data.spotify_track.snowblind.id,
  ]
}

data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
}
data "spotify_track" "blackwater" {
  url = "https://open.spotify.com/track/4lE6N1E0L8CssgKEUCgdbA"
}
data "spotify_track" "snowblind" {
  url = "https://open.spotify.com/track/7FCG2wIYG1XvGRUMACC2cD"
}
