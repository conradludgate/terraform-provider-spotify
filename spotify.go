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
	// req, err := http.NewRequest("POST", d.Get("auth_server").(string), nil)
	// if err != nil {
	// 	return nil, err
	// }
	// req.SetBasicAuth("SpotifyAuthProxy", d.Get("api_key").(string))
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("%s", string(body))
	// }

	// tokenData := struct {
	// 	AccessToken  string `json:"access_token"`
	// 	RefreshToken string `json:"refresh_token"`
	// 	ExpiresIn    int    `json:"expires_in"`
	// 	TokenType    string `json:"token_type"`
	// }{}

	// if err := json.Unmarshal(body, &tokenData); err != nil {
	// 	return nil, err
	// }

	// token := &oauth2.Token{
	// 	AccessToken:  tokenData.AccessToken,
	// 	RefreshToken: tokenData.RefreshToken,
	// 	TokenType:    tokenData.TokenType,
	// 	Expiry:       time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
	// }

	// cnf := &oauth2.Config{
	// 	// ClientID: d.Get("client_id").(string),
	// 	ClientID:     "SpotifyAuthProxy",
	// 	ClientSecret: d.Get("api_key").(string),
	// 	Endpoint: oauth2.Endpoint{
	// 		TokenURL:  d.Get("auth_server").(string),
	// 		AuthStyle: oauth2.AuthStyleInHeader,
	// 	},
	// }

	transport := &transport{
		APIKey: d.Get("api_key").(string),
		Server: d.Get("auth_server").(string),
	}
	transport.getToken()

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
		t.getToken()
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
