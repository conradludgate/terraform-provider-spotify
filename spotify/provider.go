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
				Default:     "http://localhost:27228",
				Description: "Oauth2 Proxy URL",
			},
			"token_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Oauth2 Proxy token ID",
				DefaultFunc: func() (interface{}, error) {
					return "terraform", nil
				},
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Oauth2 Proxy username",
				DefaultFunc: func() (interface{}, error) {
					return "SpotifyAuthProxy", nil
				},
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
