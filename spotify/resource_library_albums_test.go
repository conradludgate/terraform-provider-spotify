package spotify_test

import (
	"fmt"
	"testing"

	spotifyApi "github.com/conradludgate/spotify/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestSpotify_Resource_LibraryAlbums(t *testing.T) {
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

					resource "spotify_library_albums" "my_albums" {
						albums = ["album-1", "album-2"]
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "id", "library"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.#", "2"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.0", "album-1"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.1", "album-2"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("PUT", "https://api.spotify.com/v1/me/albums?ids=album-1,album-2",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/me/albums",
						RespondWith(
							JSON(savedAlbumPage("album-1", "album-2")),
							VerifyBearer(accessToken),
						),
					)
				},
			},
			{
				Config: fmt.Sprintf(`
					provider "spotify" {
						api_key = "%s"
					}

					resource "spotify_library_albums" "my_albums" {
						albums = ["album-1", "album-3"]
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "id", "library"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.#", "2"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.0", "album-1"),
					resource.TestCheckResourceAttr("spotify_library_albums.my_albums", "albums.1", "album-3"),
				),
				PreConfig: func() {
					httpmock.RegisterResponder("PUT", "https://api.spotify.com/v1/me/albums?ids=album-3",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("DELETE", "https://api.spotify.com/v1/me/albums?ids=album-2",
						RespondWith(
							JSON(nil),
							VerifyBearer(accessToken),
						).Once(),
					)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/me/albums",
						RespondWith(
							JSON(savedAlbumPage("album-1", "album-3")),
							VerifyBearer(accessToken),
						),
					)
				},
			},
		},
	})
}

func savedAlbumPage(albums ...string) spotifyApi.SavedAlbumPage {
	savedAlbums := make([]spotifyApi.SavedAlbum, len(albums))
	for i, album := range albums {
		savedAlbums[i].ID = spotifyApi.ID(album)
	}
	return spotifyApi.SavedAlbumPage{Albums: savedAlbums}
}
