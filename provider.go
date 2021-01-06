package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider for spotify
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_code": {
				Type:         schema.TypeString,
				Required:     true,
				ExactlyOneOf: []string{"client_secret", "code_verifier"},
				Description:  "An authorization code required to exchange into an access token",
			},
			"client_id": {
				Type:         schema.TypeString,
				RequiredWith: []string{"redirect_uri", "client_secret"},
				Default:      "956aed6fce0c49ebb0eb1d050d9223ed",
				Description:  "The client ID used to exchange the auth code into an access token",
			},
			"redirect_uri": {
				Type:         schema.TypeString,
				RequiredWith: []string{"client_id", "client_secret"},
				Default:      "http://localhost:27228/spotify_callback",
				Description:  "The URI that spotify redirects to when authorizing and creating the auth_code",
			},
			"code_verifier": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_secret"},
				Description:   "The code verifier value to exchange the auth code into an access token. See https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow-with-proof-key-for-code-exchange-pkce",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"code_verifier"},
				Description:   "The client secret used to convert the auth code into an access token. See https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow",
			},
		},
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
