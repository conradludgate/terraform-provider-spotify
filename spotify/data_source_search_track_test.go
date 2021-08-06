package spotify_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
	spotifyApi "github.com/zmb3/spotify"
)

func TestSpotify_DataSource_SearchTrack(t *testing.T) {
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

					data "spotify_search_track" "delta_heavy" {
						artist   = "Delta Heavy"
						explicit = false
					}
				`, apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.#", "3"),

					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.0.id", "track-1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.0.name", "White Flag"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.0.artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.0.artists.0", "artist-1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.0.album", "album-1"),

					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.1.id", "track-2"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.1.name", "Kaleidoscope"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.1.artists.#", "1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.1.artists.0", "artist-1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.1.album", "album-2"),

					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.id", "track-4"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.name", "Revenge"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.artists.#", "2"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.artists.0", "artist-1"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.artists.1", "artist-3"),
					resource.TestCheckResourceAttr("data.spotify_search_track.delta_heavy", "tracks.2.album", "album-4"),
				),
				PreConfig: func() {
					RegisterAuthResponse(apiKey, accessToken)

					httpmock.RegisterResponder("GET", "https://api.spotify.com/v1/search?limit=10&q=artist%3ADelta+Heavy&type=track",
						RespondWith(
							JSON(spotifyApi.SearchResult{
								Tracks: fullTrackPage([]track{
									{
										id:      "track-1",
										name:    "White Flag",
										artists: []string{"artist-1"},
										album:   "album-1",
									},
									{
										id:      "track-2",
										name:    "Kaleidoscope",
										artists: []string{"artist-1"},
										album:   "album-2",
									},
									{
										id:       "track-3",
										name:     "Anarchy",
										artists:  []string{"artist-1", "artist-2"},
										album:    "album-3",
										explicit: true,
									},
									{
										id:      "track-4",
										name:    "Revenge",
										artists: []string{"artist-1", "artist-3"},
										album:   "album-4",
									},
								}),
							}),
							VerifyBearer(accessToken),
						),
					)
				},
			},
		},
	})
}

type track struct {
	id       string
	name     string
	artists  []string
	album    string
	explicit bool
}

func fullTrackPage(tracks []track) *spotifyApi.FullTrackPage {
	fullTracks := make([]spotifyApi.FullTrack, len(tracks))
	for i, track := range tracks {
		fullTracks[i].SimpleTrack.ID = spotifyApi.ID(track.id)
		fullTracks[i].SimpleTrack.Name = track.name
		fullTracks[i].Album.ID = spotifyApi.ID(track.album)
		fullTracks[i].Explicit = track.explicit
		for _, artist := range track.artists {
			fullTracks[i].Artists = append(fullTracks[i].Artists, spotifyApi.SimpleArtist{
				ID: spotifyApi.ID(artist),
			})
		}
	}
	return &spotifyApi.FullTrackPage{Tracks: fullTracks}
}
