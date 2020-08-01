package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var spotifyClient *spotify.Client

// ClientConfigurer for spotify API access
func ClientConfigurer(_ *schema.ResourceData) (interface{}, error) {
	if spotifyClient != nil {
		return spotifyClient, nil
	}

	accessToken := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("SPOTIFY_ACCESS_TOKEN must be set with a valid access token")
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	httpClient := &http.Client{
		Transport: transport{
			token,
			&http.Client{},
		},
	}

	client := spotify.NewClient(httpClient)
	spotifyClient = &client
	return spotifyClient, nil
}

type transport struct {
	token *oauth2.Token
	base  *http.Client
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.token.SetAuthHeader(req)
	return t.base.Do(req)
}
