resource "spotify_playlist" "playlist" {
    name = "My playlist"
    description = "My playlist is so awesome"
    public = false

    tracks = [
        data.spotify_search_track.overkill.id,
        data.spotify_search_track.blackwater.id,
    ]
}

data "spotify_search_track" "overkill" {
    name = "Overkill"

    artists = [
        "RIOT",
    ]
}

data "spotify_search_track" "blackwater" {
    name = "Blackwater"

    artists = [
        "RIOT",
    ]
}