package spotify_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
	spotifyApi "github.com/zmb3/spotify"
)

func TestSpotify_Resource_Playlist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := "some-api-key"
	accessToken := "some-access-token"

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "spotify" {
						api_key = "%s"
					}

					resource "spotify_playlist" "playlist" {
						name = "My Playlist"
						description = "A test playlist"

						tracks = ["track-1", "track-2"]
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "id", "spotify-playlist-1"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "name", "My Playlist"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "description", "A test playlist"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "snapshot_id", "snapshot1"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.#", "2"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.0", "track-1"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.1", "track-2"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "public", "true"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/me", RespondWith(
						JSON(spotifyApi.PrivateUser{
							User: spotifyApi.User{
								ID: "user-1",
							},
						}),
						VerifyBearer(accessToken),
					))

					httpmock.RegisterResponder("POST", "https://api.spotify.com/v1/users/user-1/playlists",
						RespondWith(
							JSON(spotifyApi.FullPlaylist{
								SimplePlaylist: spotifyApi.SimplePlaylist{
									ID: spotifyApi.ID("spotify-playlist-1"),
								},
							}),
							VerifyBearer(accessToken),
							VerifyJSONBody(object{
								"name":        "My Playlist",
								"description": "A test playlist",
								"public":      true,
							}),
						).Once(),
					)

					httpmock.RegisterResponder("POST", "https://api.spotify.com/v1/playlists/spotify-playlist-1/tracks",
						RespondWith(
							JSON(object{"snapshot_id": "snapshot1"}),
							VerifyBearer(accessToken),
							VerifyJSONBody(object{
								"uris": array{
									"spotify:track:track-1",
									"spotify:track:track-2",
								},
							}),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/playlists/spotify-playlist-1",
						RespondWith(
							JSON(spotifyApi.FullPlaylist{
								SimplePlaylist: spotifyApi.SimplePlaylist{
									ID:         spotifyApi.ID("spotify-playlist-1"),
									Name:       "My Playlist",
									IsPublic:   true,
									SnapshotID: "snapshot1",
								},
								Description: "A test playlist",
							}),
							VerifyBearer(accessToken),
						),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/playlists/spotify-playlist-1/tracks",
						RespondWith(
							JSON(playlistTrackPage("track-1", "track-2")),
							VerifyBearer(accessToken),
						),
					)
				},
			},
			{
				Config: `
					provider "spotify" {
						api_key = "some-api-key"
					}

					resource "spotify_playlist" "playlist" {
						name = "My New Playlist"
						description = "A test playlist"

						tracks = ["track-1", "track-3"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "id", "spotify-playlist-1"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "name", "My New Playlist"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "description", "A test playlist"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "snapshot_id", "snapshot3"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.#", "2"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.0", "track-1"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "tracks.1", "track-3"),
					resource.TestCheckResourceAttr("spotify_playlist.playlist", "public", "true"),
				),
				PreConfig: func() {
					httpmock.RegisterResponder("PUT", "https://api.spotify.com/v1/playlists/spotify-playlist-1",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
							VerifyJSONBody(object{
								"name":        "My New Playlist",
								"description": "A test playlist",
								"public":      true,
							}),
						).Once(),
					)

					httpmock.RegisterResponder("POST", "https://api.spotify.com/v1/playlists/spotify-playlist-1/tracks",
						RespondWith(
							JSON(object{"snapshot_id": "snapshot2"}),
							VerifyBearer(accessToken),
							VerifyJSONBody(object{
								"uris": array{
									"spotify:track:track-3",
								},
							}),
						).Once(),
					)

					httpmock.RegisterResponder("DELETE", "https://api.spotify.com/v1/playlists/spotify-playlist-1/tracks",
						RespondWith(
							JSON(object{"snapshot_id": "snapshot3"}),
							VerifyBearer(accessToken),
							VerifyJSONBody(
								object{
									"tracks": array{
										object{
											"uri": "spotify:track:track-2",
										},
									},
								},
							),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/playlists/spotify-playlist-1",
						RespondWith(
							JSON(spotifyApi.FullPlaylist{
								SimplePlaylist: spotifyApi.SimplePlaylist{
									ID:         spotifyApi.ID("spotify-playlist-1"),
									Name:       "My New Playlist",
									IsPublic:   true,
									SnapshotID: "snapshot3",
								},
								Description: "A test playlist",
							}),
							VerifyBearer(accessToken),
						),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/playlists/spotify-playlist-1/tracks",
						RespondWith(
							JSON(playlistTrackPage("track-1", "track-3")),
							VerifyBearer(accessToken),
						),
					)
				},
			},
		},
	})
}

func playlistTrackPage(tracks ...string) spotifyApi.PlaylistTrackPage {
	playlistTracks := make([]spotifyApi.PlaylistTrack, len(tracks))
	for i, track := range tracks {
		playlistTracks[i].Track.ID = spotifyApi.ID(track)
	}
	return spotifyApi.PlaylistTrackPage{Tracks: playlistTracks}
}
