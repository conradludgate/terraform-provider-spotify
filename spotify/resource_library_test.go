package spotify_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
	spotifyApi "github.com/zmb3/spotify"
)

func TestSpotify_Resource_Library(t *testing.T) {
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

					resource "spotify_library" "my_library" {
						tracks = ["track-1", "track-2"]
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_library.my_library", "id", "library"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.#", "2"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.0", "track-1"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.1", "track-2"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("PUT", "https://api.spotify.com/v1/me/tracks?ids=track-1,track-2",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/me/tracks",
						RespondWith(
							JSON(savedTrackPage("track-1", "track-2")),
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

					resource "spotify_library" "my_library" {
						tracks = ["track-1", "track-3"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_library.my_library", "id", "library"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.#", "2"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.0", "track-1"),
					resource.TestCheckResourceAttr("spotify_library.my_library", "tracks.1", "track-3"),
				),
				PreConfig: func() {
					httpmock.RegisterResponder("PUT", "https://api.spotify.com/v1/me/tracks?ids=track-3",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("DELETE", "https://api.spotify.com/v1/me/tracks?ids=track-2",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/me/tracks",
						RespondWith(
							JSON(savedTrackPage("track-1", "track-3")),
							VerifyBearer(accessToken),
						),
					)
				},
			},
		},
	})
}

func savedTrackPage(tracks ...string) spotifyApi.SavedTrackPage {
	savedTracks := make([]spotifyApi.SavedTrack, len(tracks))
	for i, track := range tracks {
		savedTracks[i].FullTrack.ID = spotifyApi.ID(track)
	}
	return spotifyApi.SavedTrackPage{Tracks: savedTracks}
}
