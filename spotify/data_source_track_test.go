package spotify_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
	spotifyApi "github.com/zmb3/spotify"
)

func TestSpotify_DataSource_Track(t *testing.T) {
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

					data "spotify_track" "track-1" {
						spotify_id = "4lE6N1E0L8CssgKEUCgdbA"
					}

					data "spotify_track" "track-2" {
						url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.spotify_track.track-1", "id", "4lE6N1E0L8CssgKEUCgdbA"),
					resource.TestCheckResourceAttr("data.spotify_track.track-1", "name", "Blackwater"),
					resource.TestCheckResourceAttr("data.spotify_track.track-1", "album", "1AUS845POFhV3oDytPImEZ"),
					resource.TestCheckResourceAttr("data.spotify_track.track-1", "artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_track.track-1", "artists.0", "0qPGd8tOMHlFZt8EA1uLFY"),

					resource.TestCheckResourceAttr("data.spotify_track.track-2", "id", "4XdaaDFE881SlIaz31pTAG"),
					resource.TestCheckResourceAttr("data.spotify_track.track-2", "name", "Overkill"),
					resource.TestCheckResourceAttr("data.spotify_track.track-2", "album", "64ey3KHg3uepidKmJrb4ka"),
					resource.TestCheckResourceAttr("data.spotify_track.track-2", "artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_track.track-2", "artists.0", "0qPGd8tOMHlFZt8EA1uLFY"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/tracks/4lE6N1E0L8CssgKEUCgdbA",
						RespondWith(
							JSON(spotifyApi.FullTrack{
								SimpleTrack: spotifyApi.SimpleTrack{
									ID:   "4lE6N1E0L8CssgKEUCgdbA",
									Name: "Blackwater",
									Artists: []spotifyApi.SimpleArtist{
										{
											ID: "0qPGd8tOMHlFZt8EA1uLFY",
										},
									},
								},
								Album: spotifyApi.SimpleAlbum{
									ID: "1AUS845POFhV3oDytPImEZ",
								},
							}),
							VerifyBearer(accessToken),
						),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/tracks/4XdaaDFE881SlIaz31pTAG",
						RespondWith(
							JSON(spotifyApi.FullTrack{
								SimpleTrack: spotifyApi.SimpleTrack{
									ID:   "4XdaaDFE881SlIaz31pTAG",
									Name: "Overkill",
									Artists: []spotifyApi.SimpleArtist{
										{
											ID: "0qPGd8tOMHlFZt8EA1uLFY",
										},
									},
								},
								Album: spotifyApi.SimpleAlbum{
									ID: "64ey3KHg3uepidKmJrb4ka",
								},
							}),
							VerifyBearer(accessToken),
						),
					)
				},
			},
		},
	})
}
