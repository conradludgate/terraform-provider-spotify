package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// ClientConfigurer for spotify API access
func ClientConfigurer(d *schema.ResourceData) (interface{}, error) {
	transport := &transport{
		APIKey: d.Get("api_key").(string),
		Server: d.Get("auth_server").(string),
	}

	client := spotify.NewClient(&http.Client{
		Transport: transport,
	})
	return &client, nil
}

type transport struct {
	APIKey string
	Server string
	Base   http.RoundTripper
	token  *oauth2.Token
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.token.Valid() {
		if err := t.getToken(); err != nil {
			return nil, err
		}
	}

	t.token.SetAuthHeader(req)

	return t.base().RoundTrip(req)
}

func (t *transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func (t *transport) getToken() error {
	req, err := http.NewRequest("GET", t.Server, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("SpotifyAuthProxy", t.APIKey)
	resp, err := t.base().RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	t.token = &oauth2.Token{}
	if err := json.Unmarshal(body, t.token); err != nil {
		return err
	}

	if t.token.Valid() {
		return errors.New("could not get a valid token")
	}

	return nil
}
