package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider for spotify
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:27228/api/token",
				Description: "Spotify auth proxy URL",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Spotify auth proxy API Key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spotify_playlist": resourcePlaylist(),
			"spotify_library":  resourceLibrary(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spotify_search_track": dataSourceSearchTrack(),
			"spotify_track":        dataSourceTrack(),
		},
		ConfigureFunc: ClientConfigurer,
	}
}
