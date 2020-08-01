package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider for spotify
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"spotify_playlist": resourcePlaylist(),
			"spotify_library":  resourceLibrary(),
			// "spotify_playlist_image": resourcePlaylistImage(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spotify_search_track": dataSourceSearchTrack(),
		},
		ConfigureFunc: ClientConfigurer,
	}
}
