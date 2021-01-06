package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	spotifyOauth "golang.org/x/oauth2/spotify"
)

var spotifyClient *spotify.Client

// ClientConfigurer for spotify API access
func ClientConfigurer(d *schema.ResourceData) (interface{}, error) {
	if spotifyClient != nil {
		return spotifyClient, nil
	}

	authCode := d.Get("auth_code").(string)
	clientID := d.Get("client_id").(string)
	redirectURI := d.Get("redirect_uri").(string)

	ctx := context.Background()
	cnf := oauth2.Config{
		ClientID:    clientID,
		Endpoint:    spotifyOauth.Endpoint,
		RedirectURL: redirectURI,
	}
	options := []oauth2.AuthCodeOption{}

	if codeVerifier, ok := d.GetOk("code_verifier"); ok {
		options = append(options, oauth2.SetAuthURLParam("code_verifier", codeVerifier.(string)))
	} else if clientSecret, ok := d.GetOk("client_secret"); ok {
		cnf.ClientSecret = clientSecret.(string)
	}

	token, err := cnf.Exchange(ctx, authCode, options...)
	if err != nil {
		return nil, fmt.Errorf("Could not exchange auth code: %w", err)
	}

	client := spotify.NewClient(cnf.Client(ctx, token))
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
