package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for spotify
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:27228",
				Description: "Oauth2 Proxy URL",
			},
			"token_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "terraform",
				Description: "Oauth2 Proxy token ID",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SpotifyAuthProxy",
				Description: "Oauth2 Proxy username",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Oauth2 Proxy API Key",
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
