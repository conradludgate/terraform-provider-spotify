resource "spotify_library_albums" "my_album" {
  albums = [
    data.spotify_album.only_in_dreams.id,
    data.spotify_album.the_promised_land.id,
  ]
}

data "spotify_album" "only_in_dreams" {
  spotify_id = "35axN2yrxRiycF2pA8mZaB"
}

data "spotify_album" "the_promised_land" {
  url = "https://open.spotify.com/album/3nRnJkUJYFfxcOGgU6LNci"
}
