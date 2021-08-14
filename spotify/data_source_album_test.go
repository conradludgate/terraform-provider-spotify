package spotify_test

import (
	"fmt"
	"testing"

	spotifyApi "github.com/conradludgate/spotify/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestSpotify_DataSource_Album(t *testing.T) {
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

					data "spotify_album" "album-1" {
						spotify_id = "35axN2yrxRiycF2pA8mZaB"
					}

					data "spotify_album" "album-2" {
						url = "https://open.spotify.com/album/3nRnJkUJYFfxcOGgU6LNci"
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.spotify_album.album-1", "id", "35axN2yrxRiycF2pA8mZaB"),
					resource.TestCheckResourceAttr("data.spotify_album.album-1", "name", "Only in Dreams"),
					resource.TestCheckResourceAttr("data.spotify_album.album-1", "artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_album.album-1", "artists.0", "7GvVTb8yFV0ZrdI30Qce6T"),

					resource.TestCheckResourceAttr("data.spotify_album.album-2", "id", "3nRnJkUJYFfxcOGgU6LNci"),
					resource.TestCheckResourceAttr("data.spotify_album.album-2", "name", "The Promised Land"),
					resource.TestCheckResourceAttr("data.spotify_album.album-2", "artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_album.album-2", "artists.0", "4UNnRb4LN2hGtbtMfPzMhg"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/albums/35axN2yrxRiycF2pA8mZaB",
						RespondWith(
							JSON(spotifyApi.FullAlbum{
								SimpleAlbum: spotifyApi.SimpleAlbum{
									ID:   "35axN2yrxRiycF2pA8mZaB",
									Name: "Only in Dreams",
									Artists: []spotifyApi.SimpleArtist{
										{
											ID: "7GvVTb8yFV0ZrdI30Qce6T",
										},
									},
								},
							}),
							VerifyBearer(accessToken),
						),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/albums/3nRnJkUJYFfxcOGgU6LNci",
						RespondWith(
							JSON(spotifyApi.FullAlbum{
								SimpleAlbum: spotifyApi.SimpleAlbum{
									ID:   "3nRnJkUJYFfxcOGgU6LNci",
									Name: "The Promised Land",
									Artists: []spotifyApi.SimpleArtist{
										{
											ID: "4UNnRb4LN2hGtbtMfPzMhg",
										},
									},
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
